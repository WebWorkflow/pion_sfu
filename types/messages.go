package types

import "github.com/pion/webrtc/v3"

type JOIN_ROOM struct {
	self_id string
	room_id string
}

type LEAVE_ROOM struct {
	self_id string
	room_id string
}

type OFFER struct {
	offer   webrtc.SessionDescription
	self_id string
	room_id string
}

type ANSWER struct {
	typ string
	sdp string
}

type CANDIDATE struct {
	self_id   string
	room_id   string
	candidate webrtc.ICECandidateInit
}

type WsMessage struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}
