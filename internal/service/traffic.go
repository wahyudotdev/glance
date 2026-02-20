// Package service implements the core business logic.
package service

import (
	"glance/internal/interceptor"
	"glance/internal/model"
)

// TrafficService defines the interface for managing captured network traffic.
type TrafficService interface {
	GetPage(offset, limit int) ([]*model.TrafficEntry, int)
	Clear()
}

type trafficService struct {
	store *interceptor.TrafficStore
}

// NewTrafficService creates a new TrafficService.
func NewTrafficService(store *interceptor.TrafficStore) TrafficService {
	return &trafficService{store: store}
}

func (s *trafficService) GetPage(offset, limit int) ([]*model.TrafficEntry, int) {
	return s.store.GetPage(offset, limit)
}

func (s *trafficService) Clear() {
	s.store.ClearEntries()
}
