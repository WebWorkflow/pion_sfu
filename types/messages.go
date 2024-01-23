package types

import "github.com/pion/webrtc/v3"

type JOIN_ROOM struct {
	self_id string
	room_id string
}

type LEFT_ROOM struct {
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
	candidate webrtc.ICECandidate
}
