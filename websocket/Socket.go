package websocket

import (
	"fmt"

	"net/http"
	"pion_sfu/types"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},

	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartServer() *Wserver {
	server := Wserver{
		make(map[*websocket.Conn]bool),
	}

	http.HandleFunc("/", server.wsInit)
	go http.ListenAndServe(":8080", nil)

	return &server
}

// websockets listener
func (ws *Wserver) wsInit(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	coordinator := types.NewCoordinator()
	// TODO add socket to peer
	defer conn.Close()

	fmt.Printf("Client connected")

	if err != nil {
		fmt.Printf(" with error %s", err)
		return
	}

	fmt.Println(" successfully")

	message := &WsMessage{}

	for {
		conn.ReadJSON(&message) //deserialization doesn't work on that method
		coordinator.ObtainEvent(message)

	}
}

func (ws *Wserver) answerToPeer(message string, conn *websocket.Conn) {
	conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (ws *Wserver) broadcastJSON(message *WsMessage, conn *websocket.Conn) {
	for allconn, _ := range ws.clients {
		if conn == allconn {
			continue
		} else {
			allconn.WriteJSON(&message)
		}
	}
}

type Wserver struct {
	clients map[*websocket.Conn]bool
}

func newWSServer() *Wserver {
	return &Wserver{}
}

type WsMessage struct {
	event string
	data  any
}

func NewMessage(evt string, data any) *WsMessage {
	return &WsMessage{event: evt, data: data}
}
