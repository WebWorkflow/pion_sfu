package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// var clients map[*websocket.Conn]bool
var upgrader = websocket.Upgrader{
	//Hey CORS, fuck u
	CheckOrigin: func(r *http.Request) bool {
		return true
	},

	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// websockets listener
func wsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	defer conn.Close()

	fmt.Printf("Client connected")

	if err != nil {
		fmt.Printf(" with error %s", err)
		return
	}

	//clientID:=conn.clientID

	fmt.Println(" successfully")

	for {
		mt, message, err := conn.ReadMessage() //message type int, byte[], err

		if err != nil || mt == websocket.CloseMessage {
			break
		}

		switch string(message) {
		case "offer":
			go answerToPeer(conn, "That's ur answer")
			break

		case "answer":
			go answerToPeer(conn, "Connection established")
			break

		}

	}

}

func answerToPeer(conn *websocket.Conn, message string) {
	conn.WriteMessage(websocket.TextMessage, []byte(message))
}

// func broadcast(message []byte) {
// 	for conn := range clients {
// 		conn.WriteMessage(websocket.TextMessage, message)
// 	}
// }

type WsMessage struct {
	event string
	data  string
}
