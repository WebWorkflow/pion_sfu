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

func (c *Coordinator) ObtainEvent(event string, data any) {
	switch event {
	case "join":
		go func() {
			// if room_id is not mapped => create a new room
			// insert new user to the room
		}()
	case "leave":
		go func() {
			// drop user from the room
			// if room is empty => drop a room too
		}()
	case "add-peer":
		go func() {
			// Initiate a new peer connection
			// If we have new tracks => add tracks to room.tracks
		}()
	case "remove-peer":
		go func() {
			// remove PeerConnection from the Peer struct (peer.connection)
			// if peer.streams is not empty => remove tracks
		}()
	case "add-track":
		go func() {
			// add tracks from peer.streams to room.tracks with convertation to TrackLocalStaticRTP
		}()
	case "remove-track":
		go func() {
			// remove tracks from room.tracks and from peer.streams
		}()
	}
}
