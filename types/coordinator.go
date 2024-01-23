package types

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type Lobby interface {
	CreateRoom(id string)
	RemoveRoom(id string)
	AddUserToRoom(self_id string, room_id string, socket *websocket.Conn)
	RemoveUserFromRoom(self_id string, room_id string, socket *websocket.Conn)
}

type Coordinator struct {
	sessioins map[string]*Room
}

func NewCoordinator() *Coordinator {
	return &Coordinator{sessioins: map[string]*Room{}}
}

func (coordintor *Coordinator) ShowSessions() map[string]*Room {
	return coordintor.sessioins
}

func (coordinator *Coordinator) CreateRoom(id string) {
	coordinator.sessioins[id] = NewRoom(id)
}

func (coordinator *Coordinator) RemoveRoom(id string) {
	delete(coordinator.sessioins, id)
}

func (coordinator *Coordinator) AddUserToRoom(self_id string, room_id string, socket *websocket.Conn) {
	if _, ok := coordinator.sessioins[room_id]; !ok {
		coordinator.CreateRoom(room_id)
	}
	if room, ok := coordinator.sessioins[room_id]; ok {
		// Add Peer to Room
		room.AddPeer(newPeer(self_id))
		if peer, ok := room.peers[self_id]; ok {
			// Set socket connection to Peer
			peer.SetSocket(socket)

			// Create Peer Connection
			conn, err := webrtc.NewPeerConnection(webrtc.Configuration{})
			if err != nil {
				fmt.Println("Failed to establish peer connection")
			}

			peer.SetPeerConnection(conn)

			// TODO Do we need this ?
			defer peer.connection.Close()

			// Accept one audio and one video track incoming
			for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
				if _, err := peer.connection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
					Direction: webrtc.RTPTransceiverDirectionRecvonly,
				}); err != nil {
					log.Print(err)
					return
				}
			}

			// If PeerConnection is closed remove it from global list
			peer.connection.OnConnectionStateChange(func(p webrtc.PeerConnectionState) {
				switch p {
				case webrtc.PeerConnectionStateFailed:
					if err := peer.connection.Close(); err != nil {
						log.Print(err)
					}
				case webrtc.PeerConnectionStateClosed:
					room.Signal()
				default:
				}
			})

			// When peer connection is getting the ICE -> send ICE to client
			peer.connection.OnICECandidate(func(i *webrtc.ICECandidate) {
				if i == nil {
					return
				}

				candidateString, err := json.Marshal(i.ToJSON())
				if err != nil {
					fmt.Println(err)
					return
				}

				room.SendICE(candidateString, peer.id)
			})

			peer.connection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
				// Create a track to fan out our incoming video to all peers
				trackLocal := room.AddTrack(t)
				defer room.RemoveTrack(trackLocal)

				buf := make([]byte, 1500)
				for {
					i, _, err := t.Read(buf)
					if err != nil {
						return
					}

					if _, err = trackLocal.Write(buf[:i]); err != nil {
						return
					}
				}
			})
		}

	}
}

func (coordinator *Coordinator) RemoveUserFromRoom(self_id string, room_id string) {
	if room, ok := coordinator.sessioins[room_id]; ok {
		if _, ok := room.peers[self_id]; ok {
			delete(room.peers, self_id)
		}
	}
}

func (coordinator *Coordinator) ObtainEvent(message []byte, socket *websocket.Conn) {
	wsMessage := WsMessage{}
	err := json.Unmarshal(message, &wsMessage)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch wsMessage.Event {
	case "joinRoom":
		go func() {
			data, ok := wsMessage.Data.(JOIN_ROOM)
			if !ok {
				fmt.Println("Conversion failed")
				return
			}
			coordinator.AddUserToRoom(data.self_id, data.room_id, socket)
		}()
	case "leaveRoom":
		go func() {
			data, ok := wsMessage.Data.(LEFT_ROOM)
			if !ok {
				fmt.Println("Conversion failed")
				return
			}
			coordinator.RemoveUserFromRoom(data.self_id, data.room_id)
		}()
	case "offer":
		go func() {
			data, ok := wsMessage.Data.(OFFER)
			if !ok {
				fmt.Println("Conversion failed")
				return
			}
			if room, ok := coordinator.sessioins[data.room_id]; ok {
				if peer, ok := room.peers[data.self_id]; ok {
					answer, err2 := peer.ReactOnOffer(data.offer)
					if err2 != nil {
						fmt.Println(err2)
						return
					}
					room.SendAnswer(answer, data.self_id)
				}
			}

		}()
	case "ice-candidate":
		go func() {
			data, ok := wsMessage.Data.(CANDIDATE)
			if !ok {
				fmt.Println("Conversion failed")
				return
			}
			if room, ok := coordinator.sessioins[data.room_id]; ok {
				if peer, ok := room.peers[data.self_id]; ok {
					if err := peer.connection.AddICECandidate(data.candidate.ToJSON()); err != nil {
						log.Println(err)
						return
					}
				}
			}
		}()
	}

	return
}
