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

func StartServer() *WsServer {
	server := WsServer{
		make(map[*websocket.Conn]bool),
	}

	http.HandleFunc("/", server.wsInit)
	go http.ListenAndServe(":8080", nil)

	return &server
}

// websockets listener
func (ws *WsServer) wsInit(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	coordinator := types.NewCoordinator()
	defer conn.Close()

	fmt.Printf("Client connected")

	if err != nil {
		fmt.Printf(" with error %s", err)
		return
	}

	fmt.Println(" successfully")

	message := []byte{}

	for {
		err := conn.ReadJSON(&message)
		if err != nil {
			fmt.Println(err)
			return
		}
		coordinator.ObtainEvent(message, conn)
	}
}

func (ws *WsServer) answerToPeer(message string, conn *websocket.Conn) {
	conn.WriteMessage(websocket.TextMessage, []byte(message))
}



func (ws *WsServer) broadcastJSON(message *WsMessage, conn *websocket.Conn) {
	for allconn, _ := range ws.clients {
		if conn == allconn {
			continue
		} else {
			allconn.WriteJSON(&message)
		}
	}
}

type WsServer struct {
	clients map[*websocket.Conn]bool
}



type WsMessage struct {
	Event string
	Data  any
}

func NewMessage(evt string, data any) *WsMessage {
	return &WsMessage{Event: evt, Data: data}
}
