package service

import (
	"glance/internal/interceptor"
	"glance/internal/model"
	"testing"
)

func TestTrafficService(t *testing.T) {
	repo := &mockTrafficRepo{}
	store := interceptor.NewTrafficStore(repo)
	svc := NewTrafficService(store)

	// Test GetPage
	_ = repo.Add(&model.TrafficEntry{ID: "1"})
	entries, total := svc.GetPage(0, 10)
	if total != 1 || len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", total)
	}

	// Test Clear
	svc.Clear()
	_, total = svc.GetPage(0, 10)
	if total != 0 {
		t.Errorf("Expected 0 entries after clear, got %d", total)
	}
}

func TestTrafficService_Pagination(t *testing.T) {
	repo := &mockTrafficRepo{}
	store := interceptor.NewTrafficStore(repo)
	svc := NewTrafficService(store)

	_ = repo.Add(&model.TrafficEntry{ID: "1"})
	_ = repo.Add(&model.TrafficEntry{ID: "2"})

	entries, _ := svc.GetPage(1, 1)
	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}
}
