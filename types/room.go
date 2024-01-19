package SFUtypes

import (
	"sync"

	"github.com/pion/webrtc/v3"
)

type Session interface {
	AddPeer(peer *Peer)
	RemovePeer(peer_id string)
	// Subscribe(peer *Peer)
	// Publish(peer *Peer)
	// Unsubscribe(peer *Peer)
	// Unpublish(peer *Peer)
	// GetPeers() []Peer
	AddTrack(track *webrtc.TrackRemote)
	RemoveTrack(track *webrtc.TrackRemote)
}

type Room struct {
	id     string
	mutex  sync.RWMutex
	peers  map[string]*Peer
	tracks map[string]*webrtc.TrackLocalStaticRTP
}

func NewRoom(id string) Room {
	return Room{
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
		// TODO SIGNAL
	}()

	r.peers[peer.id] = peer
}

func (r *Room) RemovePeer(peer_id string) {
	r.mutex.Lock()
	defer func() {
		r.mutex.Unlock()
		// TODO SIGNAL
	}()

	delete(r.peers, peer_id)
}

func (r *Room) AddTrack(track *webrtc.TrackRemote) {
	r.mutex.Lock()
	defer func() {
		r.mutex.Unlock()
		// TODO SIGNAL
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
		// TODO SIGNAL
	}()

	delete(r.tracks, track.ID())
}
