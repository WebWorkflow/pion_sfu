package types

import (
	"fmt"
	"github.com/pion/webrtc/v3"
	"sync"

	"github.com/gorilla/websocket"
)

type PeerInterface interface {
	SetSocket(ws_conn *websocket.Conn)
	AddRemoteTrack(track *webrtc.TrackRemote)
	RemoveRemoteTrack(track *webrtc.TrackRemote)
	SetPeerConnection(conn *webrtc.PeerConnection)
	ReactOnOffer(offer webrtc.SessionDescription)
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

func (peer *Peer) ReactOnOffer(offer_str string) (webrtc.SessionDescription, error) {
	peer.mutex.Lock()
	defer peer.mutex.Unlock()

	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offer_str,
	}
	err := peer.connection.SetRemoteDescription(offer)
	if err != nil {
		fmt.Println(err)
		return offer, err
	}
	fmt.Println("Remote Description was set for peer ", peer.id)
	answer, err := peer.connection.CreateAnswer(nil)
	_ = peer.connection.SetLocalDescription(answer)
	fmt.Println("Local Description was set for peer ", peer.id)
	if err != nil {
		return offer, err
	}
	fmt.Println("Answer was created in peer ", peer.id)
	return answer, nil

}

func (peer *Peer) ReactOnAnswer(answer_str string) error {
	peer.mutex.Lock()
	defer peer.mutex.Unlock()
	answer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  answer_str,
	}
	err := peer.connection.SetRemoteDescription(answer)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
