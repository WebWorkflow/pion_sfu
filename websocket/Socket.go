package websocket

import (
	
	"fmt"
	
	"net/http"

	"github.com/gorilla/websocket"
)




var upgrader = websocket.Upgrader{
	//Hey CORS, fuck u
	CheckOrigin: func(r *http.Request) bool {
		return true
	},

	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// websockets listener
func (ws *Wserver) wsHandler(w http.ResponseWriter, r *http.Request) {

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
		msg:= conn.ReadJSON(conn) //deserialization doesn't work on that method

		// if err != nil || mt == websocket.CloseMessage {
		// 	log.Println(err)
		// 	return
		// } else if e := json.Unmarshal(msg, &message); e != nil {
		// 	log.Println(err)
		// 	return
		// }

		switch message.event {
		case "offer":
			go func() {
				ws.broadcastJSON(&msg)
			}()
		case "answer":
			go func() {
				ws.broadcastJSON(&msg)
			}()
		case "ice-candidate":
			go func() {
				ws.broadcastJSON(&msg)
			}()
		case "join":
			go func() {}()
		case "leave":
			go func() {}()
		}
	}

}

func(ws *Wserver) answerToPeer( message string) {
	ws.myconn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (ws *Wserver) broadcastJSON( message *WsMessage){
	for allconn,_ :=range ws.clients{
		if (ws.myconn==allconn){
           continue
		} else {
			allconn.writeJSON(&message)
		}
	}
}


type Wserver struct{
	myconn *websocket.Conn
    clients map[*websocket.Conn]bool
}


type WsMessage struct {
	event string
	data  any
}

func newMessage(evt string, data any) WsMessage {
	return WsMessage{event: evt, data: data}
}



