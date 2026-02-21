package apiserver

import (
	"glance/internal/interceptor"
	"glance/internal/proxy"
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

func TestServer_Listen(_ *testing.T) {
	store := interceptor.NewTrafficStore(nil)
	p := proxy.NewProxy(":0")
	s := NewServer(store, p, ":0", nil, nil)

	go func() {
		_, _ = s.Listen(":0")
	}()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)
}
