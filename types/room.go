package types

import (
	"encoding/json"
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

	room.tracks[track.ID()] = trackLocal
	fmt.Println("Track ", track.ID(), " was added")
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
		raw, parse_err := json.Marshal(message)
		if err := peer.socket.WriteJSON(WsMessage{Event: "answer", Data: string(raw)}); err != nil && parse_err != nil {
			fmt.Println(err)
			fmt.Println(parse_err)
		}
	}
}

func (room *Room) SendOffer(message webrtc.SessionDescription, peer_id string) {
	if peer, ok := room.peers[peer_id]; ok {
		raw, parse_err := json.Marshal(message)
		if err := peer.socket.WriteJSON(WsMessage{Event: "offer", Data: string(raw)}); err != nil && parse_err != nil {
			fmt.Println(err)
			fmt.Println(parse_err)
		}
	}
}

func (room *Room) SendICE(message *webrtc.ICECandidate, peer_id string) {
	if peer, ok := room.peers[peer_id]; ok {
		fmt.Println("SENDED |ICE|: ", message.ToJSON())
		raw, parse_err := json.Marshal(message.ToJSON())
		if err := peer.socket.WriteJSON(WsMessage{Event: "candidate", Data: string(raw)}); err != nil && parse_err != nil {
			fmt.Println(err)
			fmt.Println(parse_err)
		}
	}
}

func (room *Room) BroadCast(message WsMessage, self_id string) {
	room.mutex.Lock()
	defer room.mutex.Unlock()
	for _, rec := range room.peers {
		if rec.id != self_id {
			if err := rec.socket.WriteJSON(message); err != nil {
				fmt.Println(err)
			}
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
				fmt.Println("Peer with peer_id", peer.id, "was disconnected")
				room.RemovePeer(peer.id)
				return true
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
						fmt.Println("Track", sender.Track().ID(), "was removed")
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
					if _, err := peer.connection.AddTrack(room.tracks[trackID]); err == nil {
						fmt.Println("New track are sending for peer", peer.id)
						return true
					} else {
						fmt.Println(err)
					}
				}
			}

			if peer.connection.PendingLocalDescription() != nil {
				fmt.Println(peer.connection.PendingLocalDescription())
				offer, err := peer.connection.CreateOffer(&webrtc.OfferOptions{
					OfferAnswerOptions: webrtc.OfferAnswerOptions{},
					ICERestart:         true,
				})
				if err != nil {
					fmt.Println("Error in CreateOffer: ", err)
					return true
				}
				if err = peer.connection.SetLocalDescription(offer); err != nil {
					fmt.Println("Offer: ", offer)
					fmt.Println("Cannot set LocalDescription: ", err)
					return false
				}

				offerString, err := json.Marshal(offer)
				if err != nil {
					fmt.Println("Marshalling failed: ", err)
					return true
				}

				if err = peer.socket.WriteJSON(&WsMessage{
					Event: "offer",
					Data:  string(offerString),
				}); err != nil {
					fmt.Println("Cannot write message in WsMessage: ", err)
					return true
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
			fmt.Println("Signalling finished")
			break
		}
	}
}
