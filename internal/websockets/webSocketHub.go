package websockets

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sarvochcha01/enlace-backend/internal/models"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (configure for production)
	},
}

// Client struct to manage WebSocket connections
type Client struct {
	Conn   *websocket.Conn
	Send   chan models.NotificationResponseDTO
	UserID uuid.UUID
}

type UserIDFinder interface {
	GetUserIDByFirebaseUID(firebaseUID string) (uuid.UUID, error)
}

// WebSocketHub manages active clients
type WebSocketHub struct {
	Clients    map[uuid.UUID]*Client
	Broadcast  chan models.NotificationResponseDTO
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
	authClient *auth.Client
	userFinder UserIDFinder
}

func (hub *WebSocketHub) SetUserFinder(finder UserIDFinder) {
	hub.userFinder = finder
}

// NewWebSocketHub initializes a WebSocketHub
func NewWebSocketHub(ac *auth.Client) *WebSocketHub {
	hub := &WebSocketHub{
		Clients:    make(map[uuid.UUID]*Client),
		Broadcast:  make(chan models.NotificationResponseDTO),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		authClient: ac,
	}

	return hub
}

// Run starts the WebSocketHub
func (hub *WebSocketHub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.mu.Lock()
			hub.Clients[client.UserID] = client
			hub.mu.Unlock()
		case client := <-hub.Unregister:
			hub.mu.Lock()
			if _, ok := hub.Clients[client.UserID]; ok {
				close(client.Send)
				delete(hub.Clients, client.UserID)
			}
			hub.mu.Unlock()
		case notification := <-hub.Broadcast:
			hub.mu.Lock()
			for _, client := range hub.Clients {
				client.Send <- notification
			}
			hub.mu.Unlock()
		}
	}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	token := r.URL.Query().Get("token")
	if token == "" {
		log.Printf("Missing token")
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	firebaseToken, err := h.authClient.VerifyIDToken(context.Background(), token)
	if err != nil {
		log.Printf("Invalid token: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if h.userFinder == nil {
		log.Println("UserFinder not set")
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}

	userID, err := h.userFinder.GetUserIDByFirebaseUID(firebaseToken.UID)
	if err != nil {
		log.Println("No user found")
		http.Error(w, "No user found", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	client := &Client{
		Conn:   conn,
		UserID: userID,
		Send:   make(chan models.NotificationResponseDTO),
	}

	h.Register <- client

	go client.ReadPump(h)
	go client.WritePump()
}

func (h *WebSocketHub) SendNotificationToUser(notification models.NotificationResponseDTO) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client, exists := h.Clients[notification.UserID]; exists {
		client.Send <- notification
	}
}

func (c *Client) ReadPump(hub *WebSocketHub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Client) WritePump() {
	for notification := range c.Send {
		data, err := json.Marshal(notification)
		if err != nil {
			log.Println("Failed to encode notification:", err)
			continue
		}

		if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			break
		}
	}
}
