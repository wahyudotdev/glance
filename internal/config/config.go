// Package config manages application-wide settings and their persistence.
package config

import (
	"glance/internal/model"
	"glance/internal/repository"
)

var repo repository.ConfigRepository

// Init initializes the configuration system with the provided repository.
func Init(r repository.ConfigRepository) {
	repo = r
}

// Get returns the current application configuration.
func Get() *model.Config {
	cfg, err := repo.Get()
	if err != nil {
		// Fallback to defaults if repo fails
		return &model.Config{
			ProxyAddr:  ":8000",
			APIAddr:    ":8081",
			MCPAddr:    ":8082",
			MCPEnabled: false,
		}
	}
	return cfg
}

// Save persists the provided configuration to the repository.
func Save(c *model.Config) error {
	return repo.Save(c)
}
