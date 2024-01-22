package types

import ("github.com/gorilla/websocket")
type Lobby interface {
	CreateRoom()
	RemoveRoom()
	EditRoom()
}

type Coordinator struct {
	sessioins map[string]Room
}

// Coordinator
func NewCoordinator() *Coordinator {
	return &Coordinator{sessioins: map[string]Room{}}
}

func (c *Coordinator) CreateRoom(id string) {
	c.sessioins[id] = NewRoom(id)
}

func (c *Coordinator) RemoveRoom(id string) {
	delete(c.sessioins, id)
}

func (c *Coordinator) addUserToRoom(id string,socket *websocket.Conn){
	room,exist:=c.sessioins[id]
	if exist{
		room.addUser(socket)
	} else {
        c.CreateRoom(id)
		room.addUser(socket)
	}
}

func (c *Coordinator) removeUserFromRoom(id string){
	room,exist:=c.sessioins[id]
	if exist==false{
		return
	}

	room.removeUser(id)
}
