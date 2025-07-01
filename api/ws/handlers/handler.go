package handlers

import (
	"encoding/json"
	"go-source/bootstrap"
	logger "go-source/pkg/log"
	"go-source/pkg/utils"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	UserID  string      `json:"user_id,omitempty"`
	RoomID  string      `json:"room_id,omitempty"`
}

// WSConnection represents a WebSocket connection
type WSConnection struct {
	conn    *websocket.Conn
	send    chan []byte
	userID  string
	roomID  string
	handler *WSHandler
	mu      sync.Mutex
}

// WSHandler handles WebSocket connections
type WSHandler struct {
	handlers    *bootstrap.Handlers
	connections map[*WSConnection]bool
	broadcast   chan []byte
	register    chan *WSConnection
	unregister  chan *WSConnection
	rooms       map[string]map[*WSConnection]bool
	mu          sync.RWMutex
}

// NewWSHandler creates a new WebSocket handler
func NewWSHandler(handlers *bootstrap.Handlers) *WSHandler {
	h := &WSHandler{
		handlers:    handlers,
		connections: make(map[*WSConnection]bool),
		broadcast:   make(chan []byte),
		register:    make(chan *WSConnection),
		unregister:  make(chan *WSConnection),
		rooms:       make(map[string]map[*WSConnection]bool),
	}

	go h.run()
	return h
}

// HandleWebSocket handles the main WebSocket endpoint
func (h *WSHandler) HandleWebSocket(c echo.Context) error {
	conn, err := utils.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logger.GetLogger().Error().Msgf("WebSocket upgrade failed: %v", err)
		return err
	}

	wsConn := &WSConnection{
		conn:    conn,
		send:    make(chan []byte, 256),
		handler: h,
	}

	h.register <- wsConn

	go wsConn.writePump()
	go wsConn.readPump()

	return nil
}

// HandleChatWebSocket handles chat-specific WebSocket connections
func (h *WSHandler) HandleChatWebSocket(c echo.Context) error {
	roomID := c.QueryParam("room_id")
	userID := c.QueryParam("user_id")

	if roomID == "" || userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "room_id and user_id are required"})
	}

	conn, err := utils.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logger.GetLogger().Error().Msgf("Chat WebSocket upgrade failed: %v", err)
		return err
	}

	wsConn := &WSConnection{
		conn:    conn,
		send:    make(chan []byte, 256),
		userID:  userID,
		roomID:  roomID,
		handler: h,
	}

	h.register <- wsConn

	go wsConn.writePump()
	go wsConn.readPump()

	return nil
}

// HandleNotificationWebSocket handles notification-specific WebSocket connections
func (h *WSHandler) HandleNotificationWebSocket(c echo.Context) error {
	userID := c.QueryParam("user_id")

	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id is required"})
	}

	conn, err := utils.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logger.GetLogger().Error().Msgf("Notification WebSocket upgrade failed: %v", err)
		return err
	}

	wsConn := &WSConnection{
		conn:    conn,
		send:    make(chan []byte, 256),
		userID:  userID,
		roomID:  "notifications",
		handler: h,
	}

	h.register <- wsConn

	go wsConn.writePump()
	go wsConn.readPump()

	return nil
}

// run handles the main WebSocket event loop
func (h *WSHandler) run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.connections[conn] = true

			// Add to room if specified
			if conn.roomID != "" {
				if h.rooms[conn.roomID] == nil {
					h.rooms[conn.roomID] = make(map[*WSConnection]bool)
				}
				h.rooms[conn.roomID][conn] = true
			}
			h.mu.Unlock()

			// Send welcome message
			welcomeMsg := WSMessage{
				Type:    "connected",
				Payload: "Successfully connected to WebSocket",
				UserID:  conn.userID,
			}
			h.sendToConnection(conn, welcomeMsg)

		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.connections[conn]; ok {
				delete(h.connections, conn)
				close(conn.send)
			}

			// Remove from room
			if conn.roomID != "" && h.rooms[conn.roomID] != nil {
				delete(h.rooms[conn.roomID], conn)
				if len(h.rooms[conn.roomID]) == 0 {
					delete(h.rooms, conn.roomID)
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for conn := range h.connections {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(h.connections, conn)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// sendToConnection sends a message to a specific connection
func (h *WSHandler) sendToConnection(conn *WSConnection, msg WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		logger.GetLogger().Error().Msgf("Failed to marshal message: %v", err)
		return
	}

	conn.mu.Lock()
	defer conn.mu.Unlock()

	select {
	case conn.send <- data:
	default:
		close(conn.send)
		delete(h.connections, conn)
	}
}

// BroadcastToRoom sends a message to all connections in a specific room
func (h *WSHandler) BroadcastToRoom(roomID string, msg WSMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, exists := h.rooms[roomID]; exists {
		data, err := json.Marshal(msg)
		if err != nil {
			logger.GetLogger().Error().Msgf("Failed to marshal message: %v", err)
			return
		}

		for conn := range room {
			select {
			case conn.send <- data:
			default:
				close(conn.send)
				delete(h.connections, conn)
				delete(room, conn)
			}
		}
	}
}

// BroadcastToUser sends a message to all connections of a specific user
func (h *WSHandler) BroadcastToUser(userID string, msg WSMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(msg)
	if err != nil {
		logger.GetLogger().Error().Msgf("Failed to marshal message: %v", err)
		return
	}

	for conn := range h.connections {
		if conn.userID == userID {
			select {
			case conn.send <- data:
			default:
				close(conn.send)
				delete(h.connections, conn)
			}
		}
	}
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *WSConnection) readPump() {
	defer func() {
		c.handler.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.GetLogger().Error().Msgf("WebSocket read error: %v", err)
			}
			break
		}

		// Parse the message
		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			logger.GetLogger().Error().Msgf("Failed to unmarshal message: %v", err)
			continue
		}

		// Handle different message types
		switch wsMsg.Type {
		case "chat":
			// Broadcast to room
			if c.roomID != "" {
				wsMsg.UserID = c.userID
				c.handler.BroadcastToRoom(c.roomID, wsMsg)
			}
		case "notification":
			// Handle notification
			logger.GetLogger().Info().Msgf("Notification received: %+v", wsMsg)
		case "ping":
			// Send pong response
			pongMsg := WSMessage{
				Type:    "pong",
				Payload: "pong",
			}
			c.handler.sendToConnection(c, pongMsg)
		default:
			logger.GetLogger().Info().Msgf("Unknown message type: %s", wsMsg.Type)
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *WSConnection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
