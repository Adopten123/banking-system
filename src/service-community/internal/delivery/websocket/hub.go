package websocket

import (
	"log"
	"sync"

	"github.com/Adopten123/banking-system/service-community/internal/config"
)

type Hub struct {
	mu         sync.RWMutex
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	cfg        config.WebSocketConfig
}

func NewHub(cfg config.WebSocketConfig) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		cfg:        cfg,
	}
}

func (h *Hub) Run() {
	log.Println("WebSocket Hub started...")
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("User %s connected to WebSocket. Total clients: %d", client.UserID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
				log.Printf("User %s disconnected. Total clients: %d", client.UserID, len(h.clients))
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) SendToUser(userID string, message []byte) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if ok {
		select {
		case client.Send <- message:
		default:
			log.Printf("Buffer full for user %s, dropping message", userID)
		}
	}
}
