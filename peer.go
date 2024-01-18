package main

import (
	"sync"

	"github.com/pion/webrtc/v3"
)

type Peer struct {
	id         string
	connection *webrtc.PeerConnection
	streams    map[string]*webrtc.TrackLocal
	mutex      sync.RWMutex
}
