package types

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Lobby interface {
	CreateRoom()
	RemoveRoom()
	EditRoom()
}

type Coordinator struct {
	sessioins map[string]*Room
}

func NewCoordinator() *Coordinator {
	return &Coordinator{sessioins: map[string]*Room{}}
}

func (coordinator *Coordinator) CreateRoom(id string) {
	coordinator.sessioins[id] = NewRoom(id)
}

func (coordinator *Coordinator) RemoveRoom(id string) {
	delete(coordinator.sessioins, id)
}

func (coordinator *Coordinator) addUserToRoom(id string, socket *websocket.Conn) {
	room, exist := coordinator.sessioins[id]
	peer := newPeer(socket.LocalAddr().String())
	peer.SetSocket(socket)
	if exist {
		room.AddPeer(peer)
	} else {
		coordinator.CreateRoom(id)
		room.AddPeer(peer)
	}
}

func (coordinator *Coordinator) removeUserFromRoom(id string, socketLocalAddr string) {
	room, exist := coordinator.sessioins[id]
	if !exist {
		return
	}
	room.RemovePeer(socketLocalAddr)
}

func (coordinator *Coordinator) findInRoom(LocalAddr string) (*Peer, string, error) {
	for roomID, room := range coordinator.sessioins {
		peer, exist := room.peers[LocalAddr]

		if exist {
			return peer, roomID, nil
		}
	}
	return newPeer(""), "", fmt.Errorf("Rooms don't include such user")

}

func (coordinator *Coordinator) ObtainEvent(message any) {

}
