package SFUtypes

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type PeerInterface interface {
	SetSocket(ws_conn *websocket.Conn)
	AddRemoteTrack(track *webrtc.TrackRemote)
	RemoveRemoteTrack(track *webrtc.TrackRemote)
	SetPeerConnection(conn *webrtc.PeerConnection)
}

type Peer struct {
	id         string
	connection *webrtc.PeerConnection
	streams    map[string]*webrtc.TrackRemote
	mutex      sync.RWMutex
	socket     *websocket.Conn
}

func newPeer(id string) *Peer {
	return &Peer{id: id, mutex: sync.RWMutex{}}
}

func (peer *Peer) SetPeerConnection(conn *webrtc.PeerConnection) {
	peer.mutex.Lock()
	defer peer.mutex.Unlock()
	peer.connection = conn
}

func (peer *Peer) AddRemoteTrack(track *webrtc.TrackRemote) {
	peer.mutex.Lock()
	defer peer.mutex.Unlock()
	peer.streams[track.ID()] = track
}

func (peer *Peer) RemoveRemoteTrack(track *webrtc.TrackRemote) {
	peer.mutex.Lock()
	defer peer.mutex.Unlock()
	delete(peer.streams, track.ID())
}

func (peer *Peer) SetSocket(socket *websocket.Conn) {
	peer.mutex.Lock()
	defer peer.mutex.Unlock()
	peer.socket = socket
}
