// Package repository defines the interfaces for data persistence.
package repository

import (
	"glance/internal/model"
)

// ConfigRepository defines the interface for managing application configuration.
type ConfigRepository interface {
	Get() (*model.Config, error)
	Save(cfg *model.Config) error
}

// TrafficRepository defines the interface for storing and retrieving HTTP traffic.
type TrafficRepository interface {
	Add(entry *model.TrafficEntry) error
	GetPage(offset, limit int) ([]*model.TrafficEntry, int, error)
	Clear() error
	Prune(limit int) error
	Flush() // For testing/synchronization
}

// RuleRepository defines the interface for managing interception rules.
type RuleRepository interface {
	GetAll() ([]*model.Rule, error)
	Add(rule *model.Rule) error
	Update(rule *model.Rule) error
	Delete(id string) error
}
