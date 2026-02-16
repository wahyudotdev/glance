package repository

import (
	"agent-proxy/internal/model"
)

type ConfigRepository interface {
	Get() (*model.Config, error)
	Save(cfg *model.Config) error
}

type TrafficRepository interface {
	Add(entry *model.TrafficEntry) error
	GetAll() ([]*model.TrafficEntry, error)
	Clear() error
	Prune(limit int) error
}
