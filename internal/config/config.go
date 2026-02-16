package config

import (
	"agent-proxy/internal/model"
	"agent-proxy/internal/repository"
)

var repo repository.ConfigRepository

func Init(r repository.ConfigRepository) {
	repo = r
}

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

func Save(c *model.Config) error {
	return repo.Save(c)
}
