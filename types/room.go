package types

import (
	"fmt"
	"sync"
	"time"

	"github.com/pion/webrtc/v3"
)

type Session interface {
	JoinRoom(id string)
	AddPeer(peer *Peer)
	RemovePeer(peer_id string)
	AddTrack(track *webrtc.TrackRemote)
	RemoveTrack(track *webrtc.TrackRemote)
	SendAnswer(message webrtc.SessionDescription, peer_id string)
	Signal()
}

type Room struct {
	id     string
	mutex  sync.RWMutex
	peers  map[string]*Peer
	tracks map[string]*webrtc.TrackLocalStaticRTP
}

func NewRoom(id string) *Room {
	return &Room{
		id:     id,
		mutex:  sync.RWMutex{},
		peers:  map[string]*Peer{},
		tracks: map[string]*webrtc.TrackLocalStaticRTP{},
	}
}

func (room *Room) AddPeer(peer *Peer) {
	room.mutex.Lock()
	defer func() {
		room.mutex.Unlock()
		room.Signal()
	}()

	room.peers[peer.id] = peer
}

func (room *Room) RemovePeer(peer_id string) {
	room.mutex.Lock()
	defer func() {
		room.mutex.Unlock()
		room.Signal()
	}()

	delete(room.peers, peer_id)
}

func (room *Room) AddTrack(track *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	room.mutex.Lock()
	defer func() {
		room.mutex.Unlock()
		room.Signal()
	}()
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(track.Codec().RTPCodecCapability, track.ID(), track.StreamID())
	if err != nil {
		panic(err)
	}

	room.tracks[trackLocal.ID()] = trackLocal
	return trackLocal
}

func (room *Room) RemoveTrack(track *webrtc.TrackLocalStaticRTP) {
	room.mutex.Lock()
	defer func() {
		room.mutex.Unlock()
		room.Signal()
	}()

	delete(room.tracks, track.ID())
}

func (room *Room) SendAnswer(message webrtc.SessionDescription, peer_id string) {
	if peer, ok := room.peers[peer_id]; ok {
		if err := peer.socket.WriteJSON(message); err != nil {
			fmt.Println(err)
		}
	}
}

func (room *Room) SendICE(message []byte, peer_id string) {
	if peer, ok := room.peers[peer_id]; ok {
		if err := peer.socket.WriteJSON(WsMessage{Event: "candidate", Data: message}); err != nil {
			fmt.Println(err)
		}
	}
}

func (room *Room) BroadCast(message WsMessage, self_id string) {
	room.mutex.Lock()
	defer room.mutex.Unlock()
	for _, rec := range room.peers {
		if rec.id != self_id {
			rec.socket.WriteJSON(message)
		}
	}
}

func (room *Room) JoinRoom(id string) {
	room.mutex.Lock()
	defer room.mutex.Unlock()
	room.peers[id] = newPeer(id)
}

func (room *Room) Signal() {
	room.mutex.Lock()
	defer room.mutex.Unlock()
	attemptSync := func() (again bool) {
		for _, peer := range room.peers {

			// Check if peer is already closed
			if peer.connection.ConnectionState() == webrtc.PeerConnectionStateClosed {
				room.RemovePeer(peer.id)
				return
			}

			existingSenders := map[string]bool{}
			for _, sender := range peer.connection.GetSenders() {
				if sender.Track() == nil {
					continue
				}

				existingSenders[sender.Track().ID()] = true

				// If we have a RTPSender that doesn't map to a existing track remove and signal
				if _, ok := room.tracks[sender.Track().ID()]; !ok {
					if err := peer.connection.RemoveTrack(sender); err != nil {
						return true
					}
				}
			}

			// Don't receive videos we are sending, make sure we don't have loopback
			for _, receiver := range peer.connection.GetReceivers() {
				if receiver.Track() == nil {
					continue
				}

				existingSenders[receiver.Track().ID()] = true
			}

			// Add all track we aren't sending yet to the PeerConnection
			for trackID := range room.tracks {
				if _, ok := existingSenders[trackID]; !ok {
					if _, err := peer.connection.AddTrack(room.tracks[trackID]); err != nil {
						return true
					}
				}
			}
		}
		return
	}

	for syncAttempt := 0; ; syncAttempt++ {
		if syncAttempt == 25 {
			// Release the lock and attempt a sync in 3 seconds. We might be blocking a RemoveTrack or AddTrack
			go func() {
				time.Sleep(time.Second * 3)
				room.Signal()
			}()
			return
		}

		if !attemptSync() {
			break
		}
	}
}
