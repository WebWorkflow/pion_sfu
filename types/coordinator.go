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
	sessioins map[string]Room
}

func NewCoordinator() *Coordinator {
	return &Coordinator{sessioins: map[string]Room{}}
}

func (c *Coordinator) CreateRoom(id string) {
	c.sessioins[id] = NewRoom(id)
}

func (c *Coordinator) RemoveRoom(id string) {
	delete(c.sessioins, id)
}

func (c *Coordinator) addUserToRoom(id string, socket *websocket.Conn) {
	room, exist := c.sessioins[id]
	peer := newPeer(socket.LocalAddr().String())
	peer.SetSocket(socket)
	if exist {
		room.AddPeer(peer)
	} else {
		c.CreateRoom(id)
		room.AddPeer(peer)
	}
}

func (c *Coordinator) removeUserFromRoom(id string, socketLocalAddr string) {
	room, exist := c.sessioins[id]
	if !exist {
		return
	}
	room.RemovePeer(socketLocalAddr)
}

func (c *Coordinator) findInRoom(LocalAddr string) (*Peer, string, error) {
	for roomID, room := range c.sessioins {
		peer, exist := room.peers[LocalAddr]

		if exist {
			return peer, roomID, nil
		}
	}
	return newPeer(""), "", fmt.Errorf("Rooms don't include such user")

}
