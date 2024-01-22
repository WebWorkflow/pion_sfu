package types

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	websocket2 "pion_sfu/websocket"
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

// TODO fix it
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

// TODO fix it
func (coordinator *Coordinator) removeUserFromRoom(id string, socketLocalAddr string) {
	room, exist := coordinator.sessioins[id]
	if !exist {
		return
	}
	room.RemovePeer(socketLocalAddr)
}

func (coordinator *Coordinator) ObtainEvent(message []byte) error {
	// TODO add events
	wsMessage := websocket2.WsMessage{}
	err := json.Unmarshal(message, &wsMessage)
	if err != nil {
		return fmt.Errorf("Shit")
	}

	return nil
}
