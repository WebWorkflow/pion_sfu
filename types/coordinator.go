package SFUtypes

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
