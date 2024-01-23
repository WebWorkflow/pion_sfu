package websocket

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"net/http"
	"pion_sfu/types"
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

	message := []byte{}

	for {
		err := conn.ReadJSON(&message)
		if err != nil {
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

func initReddis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()

	err := client.Set(ctx, "websocketAdr", "room1", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get(ctx, "websocketAdr").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val)
}
