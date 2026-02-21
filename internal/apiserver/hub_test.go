package apiserver

import (
	"glance/internal/model"
	"net"
	"testing"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	fiber_ws "github.com/gofiber/websocket/v2"
)

func TestHub_Logic(t *testing.T) {
	h := NewHub()
	go h.Run()

	app := fiber.New()
	app.Get("/ws", fiber_ws.New(func(c *fiber_ws.Conn) {
		h.register <- c
		defer func() { h.unregister <- c }()
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	}))

	// Start server on random port
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	addr := ln.Addr().String()
	go func() { _ = app.Listener(ln) }()
	defer func() { _ = app.Shutdown() }()

	// Connect a client
	wsURL := "ws://" + addr + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.Close()

	// Give it a moment to register
	time.Sleep(100 * time.Millisecond)

	// Test Broadcast
	h.Broadcast(&model.TrafficEntry{ID: "test"})
	h.BroadcastData([]byte("raw test"))

	// Verify message received
	_, msg, err := conn.ReadMessage()
	if err != nil || len(msg) == 0 {
		t.Errorf("Expected message, got err=%v", err)
	}

	// Test WebSocket write error by closing connection and then broadcasting
	conn2, _, _ := websocket.DefaultDialer.Dial(wsURL, nil)
	time.Sleep(100 * time.Millisecond) // wait for registration
	_ = conn2.Close()
	time.Sleep(100 * time.Millisecond) // wait for unregister or at least close
	h.BroadcastData([]byte("trigger error"))
	time.Sleep(100 * time.Millisecond)

	// Test nil safety
	h.register <- nil
	h.unregister <- nil

	// Test marshaling error
	h.Broadcast(nil)
}
