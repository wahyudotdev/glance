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
	GetPage(offset, limit int) ([]*model.TrafficEntry, int, error)
	Clear() error
	Prune(limit int) error
}
