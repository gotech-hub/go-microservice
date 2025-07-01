package utils

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader is the WebSocket upgrader configuration
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}
