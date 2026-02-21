package apiserver

import (
	"glance/internal/interceptor"
	"glance/internal/proxy"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	store := interceptor.NewTrafficStore(nil)
	p := proxy.NewProxy(":0")
	s := NewServer(store, p, ":0", nil, nil)

	if s == nil {
		t.Fatal("Expected server instance, got nil")
	}

	// Verify routes are registered without panic
	s.RegisterRoutes()
}

func TestServer_Listen(t *testing.T) {
	store := interceptor.NewTrafficStore(nil)
	p := proxy.NewProxy(":0")
	s := NewServer(store, p, ":0", nil, nil)

	// Test successful listen on random port
	// We run it in a goroutine and then shutdown quickly
	go func() {
		time.Sleep(100 * time.Millisecond)
		_ = s.app.Shutdown()
	}()

	addr, err := s.Listen(":0")
	if err != nil && err != http.ErrServerClosed {
		t.Fatalf("Listen failed: %v", err)
	}
	if addr == "" {
		t.Error("Expected address")
	}

	// Test port collision fallback
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	occupiedAddr := ln.Addr().String()
	defer func() { _ = ln.Close() }()

	s2 := NewServer(store, p, ":0", nil, nil)
	go func() {
		time.Sleep(100 * time.Millisecond)
		_ = s2.app.Shutdown()
	}()

	addr2, err2 := s2.Listen(occupiedAddr)
	if err2 != nil && err2 != http.ErrServerClosed {
		t.Fatalf("Listen fallback failed: %v", err2)
	}
	if addr2 == occupiedAddr {
		t.Error("Expected different address for fallback")
	}
}
