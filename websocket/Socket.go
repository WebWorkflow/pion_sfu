package websocket

import (
	"encoding/json"
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
		*types.NewCoordinator(),
	}
	http.HandleFunc("/", server.wsInit)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		fmt.Println(err)
	}
	return &server
}

// websockets listener
func (ws *WsServer) wsInit(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	defer conn.Close()

	fmt.Printf("Client connected")

	if err != nil {
		fmt.Printf(" with error %s", err)
		return
	}

	fmt.Println(" successfully")

	message := types.WsMessage{}

	for {
		messageType, bmessage, err := conn.ReadMessage()

		if err != nil {
			//fmt.Println("DROP2")
			//fmt.Println(messageType)
			fmt.Println(err)
			return
		}
		if messageType == websocket.CloseMessage {
			break
		}

		err = json.Unmarshal(bmessage, &message)
		if err != nil {
			fmt.Println("DROP")
			fmt.Println(message.Data)
			fmt.Println(err)
			return
		}
		ws.coordinator.ObtainEvent(message, conn)
	}
}

func (ws *WsServer) answerToPeer(message string, conn *websocket.Conn) {
	conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (ws *WsServer) broadcastJSON(message *types.WsMessage, conn *websocket.Conn) {
	for allconn, _ := range ws.clients {
		if conn == allconn {
			continue
		} else {
			allconn.WriteJSON(&message)
		}
	}
}

type WsServer struct {
	clients     map[*websocket.Conn]bool
	coordinator types.Coordinator
}
