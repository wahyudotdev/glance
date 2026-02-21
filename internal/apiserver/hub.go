package apiserver

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
	"glance/internal/model"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// Run starts the Hub and handles registration, unregistration, and broadcasting.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if client == nil {
				continue
			}
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Println("WebSocket client registered")

		case client := <-h.unregister:
			if client == nil {
				continue
			}
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				_ = client.Close()
			}
			h.mu.Unlock()
			log.Println("WebSocket client unregistered")

		case data := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Printf("WebSocket write error: %v", err)
					_ = client.Close()
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// Broadcast sends a traffic entry to all registered clients.
func (h *Hub) Broadcast(entry *model.TrafficEntry) {
	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Error marshaling traffic entry: %v", err)
		return
	}
	h.broadcast <- data
}

// BroadcastData sends raw data to all registered clients.
func (h *Hub) BroadcastData(data []byte) {
	h.broadcast <- data
}
