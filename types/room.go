package types

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/pion/webrtc/v3"

	"pion_sfu/websocket"
)

type Session interface {
	AddPeer(peer *Peer)
	RemovePeer(peer_id string)
	AddTrack(track *webrtc.TrackRemote)
	RemoveTrack(track *webrtc.TrackRemote)
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

func (r *Room) AddPeer(peer *Peer) {
	r.mutex.Lock()
	defer func() {
		r.mutex.Unlock()
		r.Signal()
	}()

	r.peers[peer.id] = peer
}

func (r *Room) RemovePeer(peer_id string) {
	r.mutex.Lock()
	defer func() {
		r.mutex.Unlock()
		r.Signal()
	}()

	delete(r.peers, peer_id)
}

func (r *Room) AddTrack(track *webrtc.TrackRemote) {
	r.mutex.Lock()
	defer func() {
		r.mutex.Unlock()
		r.Signal()
	}()
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(track.Codec().RTPCodecCapability, track.ID(), track.StreamID())
	if err != nil {
		panic(err)
	}

	r.tracks[trackLocal.ID()] = trackLocal
}

func (r *Room) RemoveTrack(track *webrtc.TrackLocalStaticRTP) {
	r.mutex.Lock()
	defer func() {
		r.mutex.Unlock()
		r.Signal()
	}()

	delete(r.tracks, track.ID())
}

func (room *Room) BroadCast(message websocket.WsMessage) {
	room.mutex.Lock()
	defer room.mutex.Unlock()
	for _, rec := range room.peers {
		rec.socket.WriteJSON(message)
	}
}

// async signaling??
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

			offer, err := peer.connection.CreateOffer(nil)
			if err != nil {
				return true
			}

			if err = peer.connection.SetLocalDescription(offer); err != nil {
				return true
			}

			offerString, err := json.Marshal(offer)
			if err != nil {
				return true
			}
			msg := websocket.NewMessage("offer", offerString)
			room.BroadCast(*msg)
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
