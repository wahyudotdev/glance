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
	GetByIDs(ids []string) ([]*model.TrafficEntry, error)
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

// ScenarioRepository defines the interface for managing recorded traffic scenarios.
type ScenarioRepository interface {
	GetAll() ([]*model.Scenario, error)
	GetByID(id string) (*model.Scenario, error)
	Add(scenario *model.Scenario) error
	Update(scenario *model.Scenario) error
	Delete(id string) error
}
