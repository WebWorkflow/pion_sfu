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

	
// websockets listener
func (ws *Wserver) wsInit(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	coordinator:=types.NewCoordinator()
	

	defer conn.Close()

	fmt.Printf("Client connected")

	if err != nil {
		fmt.Printf(" with error %s", err)
		return
	}

	fmt.Println(" successfully")

	

	message := &WsMessage{}
	
	for {
		message= conn.ReadJSON(conn) //deserialization doesn't work on that method

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
			go func() {
             coordinator.addUserToRoom(message.data,conn)
			}()
		case "leave":
			go func() {
				coordinator.removeUserFromRoom(message.data)
			}()
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
    clients map[*websocket.Conn] bool
}

func newWSServer () *Wserver{
	return &Wserver{}
}


type WsMessage struct {
	event string
	data  any
}

func newMessage(evt string, data any) *WsMessage {
	return &WsMessage{event: evt, data: data}
}




