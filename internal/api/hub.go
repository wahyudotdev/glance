package api

import (
	"encoding/json"
	"log"
	"sync"

	"agent-proxy/internal/model"
	"github.com/gofiber/websocket/v2"
)

type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan *model.TrafficEntry
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan *model.TrafficEntry),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Println("WebSocket client registered")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
			log.Println("WebSocket client unregistered")

		case entry := <-h.broadcast:
			data, err := json.Marshal(entry)
			if err != nil {
				log.Printf("Error marshaling traffic entry: %v", err)
				continue
			}

			h.mu.Lock()
			for client := range h.clients {
				if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Printf("WebSocket write error: %v", err)
					client.Close()
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) Broadcast(entry *model.TrafficEntry) {
	h.broadcast <- entry
}
