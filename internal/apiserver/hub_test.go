package apiserver

import (
	"glance/internal/model"
	"testing"
	"time"
)

func TestHub_Logic(_ *testing.T) {
	h := NewHub()
	go h.Run()

	// Exercise broadcast channels
	h.Broadcast(&model.TrafficEntry{ID: "test"})
	h.BroadcastData([]byte("test"))

	time.Sleep(50 * time.Millisecond)
}
