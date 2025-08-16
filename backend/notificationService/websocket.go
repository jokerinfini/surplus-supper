package notificationService

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// WebSocketHandler handles WebSocket connections for real-time notifications
func (s *NotificationService) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from query parameters
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id parameter", http.StatusBadRequest)
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Create client
	client := &Client{
		ID:     int(time.Now().UnixNano()), // Simple ID generation
		UserID: userID,
		Send:   make(chan []byte, 256),
		Hub:    s,
	}

	// Register client
	s.RegisterClient(client)

	// Start goroutines for reading and writing
	go client.writePump(conn)
	go client.readPump(conn)
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump(conn *websocket.Conn) {
	defer func() {
		c.Hub.UnregisterClient(c)
		conn.Close()
	}()

	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		// Handle incoming messages (e.g., mark notification as read)
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// Handle different message types
		if msgType, ok := msg["type"].(string); ok {
			switch msgType {
			case "mark_read":
				if notificationID, ok := msg["notification_id"].(float64); ok {
					_, err := c.Hub.MarkNotificationAsRead(int(notificationID))
					if err != nil {
						log.Printf("Failed to mark notification as read: %v", err)
					}
				}
			case "ping":
				// Send pong response
				response := map[string]interface{}{
					"type": "pong",
					"timestamp": time.Now().Unix(),
				}
				responseJSON, _ := json.Marshal(response)
				c.Send <- responseJSON
			}
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump(conn *websocket.Conn) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current WebSocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendNotification sends a notification to a specific client
func (s *NotificationService) SendNotification(clientID int, notification Notification) {
	s.mutex.RLock()
	client, exists := s.clients[clientID]
	s.mutex.RUnlock()

	if !exists {
		return
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return
	}

	select {
	case client.Send <- notificationJSON:
	default:
		// Client buffer is full, skip this notification
		log.Printf("Client %d buffer full, skipping notification", clientID)
	}
}

// BroadcastToAll sends a message to all connected clients
func (s *NotificationService) BroadcastToAll(message map[string]interface{}) {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal broadcast message: %v", err)
		return
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, client := range s.clients {
		select {
		case client.Send <- messageJSON:
		default:
			// Client buffer is full, skip this message
			log.Printf("Client %d buffer full, skipping broadcast", client.ID)
		}
	}
}

// GetConnectedClientsCount returns the number of connected clients
func (s *NotificationService) GetConnectedClientsCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.clients)
}

// GetConnectedUsersCount returns the number of unique connected users
func (s *NotificationService) GetConnectedUsersCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	userMap := make(map[int]bool)
	for _, client := range s.clients {
		userMap[client.UserID] = true
	}

	return len(userMap)
} 