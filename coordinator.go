package main

type Lobby interface {
	CreateRoom()
	RemoveRoom()
	EditRoom()
}

type Coordinator struct {
	sessioins map[int]Room
}

// Coordinator
func NewCoordinator() *Coordinator {
	return &Coordinator{sessioins: map[int]Room{}}
}
