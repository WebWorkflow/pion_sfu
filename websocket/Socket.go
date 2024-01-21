package websocket

import (
	"encoding/json"
	"fmt"
	"log"
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

	message := &WsMessage{}
	for {
		mt, msg, err := conn.ReadMessage() //message type int, byte[], err

		if err != nil || mt == websocket.CloseMessage {
			log.Println(err)
			return
		} else if e := json.Unmarshal(msg, &message); e != nil {
			log.Println(err)
			return
		}

		switch message.event {
		case "offer":
			go func() {}()
		case "answer":
			go func() {}()
		case "ice-candidate":
			go func() {}()
		case "join":
			go func() {}()
		case "leave":
			go func() {}()
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
	data  any
}

func newMessage(evt string, data any) WsMessage {
	return WsMessage{event: evt, data: data}
}
