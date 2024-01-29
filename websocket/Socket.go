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
	fmt.Println("Server started successfully")
	//err := http.ListenAndServe("localhost:8080", nil)
	err := http.ListenAndServe("0.0.0.0:8080", nil)
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

type WsServer struct {
	clients     map[*websocket.Conn]bool
	coordinator types.Coordinator
}
