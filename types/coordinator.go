package types

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	websocket2 "pion_sfu/websocket"
)

type Lobby interface {
	CreateRoom(id string)
	RemoveRoom(id string)
	AddUserToRoom(self_id string, room_id string, socket *websocket.Conn)
	RemoveUserFromRoom(self_id string, room_id string, socket *websocket.Conn)
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

func (coordinator *Coordinator) AddUserToRoom(self_id string, room_id string, socket *websocket.Conn) {
	if _, ok := coordinator.sessioins[room_id]; !ok {
		coordinator.CreateRoom(room_id)
	}
	if room, ok := coordinator.sessioins[room_id]; ok {
		room.AddPeer(newPeer(self_id))
		if peer, ok := room.peers[self_id]; ok {
			peer.SetSocket(socket)
		}
	}
}

func (coordinator *Coordinator) RemoveUserFromRoom(self_id string, room_id string) {
	if room, ok := coordinator.sessioins[room_id]; ok {
		if _, ok := room.peers[self_id]; ok {
			delete(room.peers, self_id)
		}
	}
}

func (coordinator *Coordinator) ObtainEvent(message []byte, socket *websocket.Conn) error {
	// TODO add events
	wsMessage := websocket2.WsMessage{}
	err := json.Unmarshal(message, &wsMessage)
	if err != nil {
		return fmt.Errorf("Shit")
	}

	switch wsMessage.Event {
	case "joinRoom":
		go func() {
			data, ok := wsMessage.Data.(JOIN_ROOM)
			if !ok {
				fmt.Println("Conversion failed")
				return
			}
			coordinator.AddUserToRoom(data.self_id, data.room_id, socket)
		}()
	case "leaveRoom":
		go func() {
			data, ok := wsMessage.Data.(LEFT_ROOM)
			if !ok {
				fmt.Println("Conversion failed")
				return
			}
			coordinator.RemoveUserFromRoom(data.self_id, data.room_id)
		}()
	case "offer":
		go func() {
			data, ok := wsMessage.Data.(OFFER)
			if !ok {
				fmt.Println("Conversion failed")
				return
			}
			if room, ok := coordinator.sessioins[data.room_id]; ok {
				if peer, ok := room.peers[data.self_id]; ok {
					answer, err2 := peer.ReactOnOffer(data.offer)
					if err2 != nil {
						fmt.Println(err2)
						return
					}
					room.SendAnswer(answer, data.self_id)
				}
			}

		}()
	case "ice-candidate":
		go func() {

		}()
	}

	return nil
}
