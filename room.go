package main

import (
	"sync"
)

type Session interface {
	AddPeer(peer *Peer)
	RemovePeer(peer_id string)
	Subscribe(peer *Peer)
	Publish(peer *Peer)
	Unsubscribe(peer *Peer)
	Unpublish(peer *Peer)
	GetPeers() []Peer
}

type Room struct {
	id    string
	mutex sync.RWMutex
	peers map[string]*Peer
}

func (r *Room) AddPeer(peer *Peer) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.peers[peer.id] = peer
}

func (r *Room) RemovePeer(peer_id string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.peers, peer_id)
}
