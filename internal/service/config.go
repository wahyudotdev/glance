// Package service implements the core business logic.
package service

import (
	"glance/internal/config"
	"glance/internal/mcp"
	"glance/internal/model"
)

// ConfigService defines the interface for application configuration and status.
type ConfigService interface {
	GetStatus(mcpServer *mcp.Server, proxyAddr string) (map[string]any, error)
	GetConfig() (*model.Config, error)
	SaveConfig(cfg *model.Config) error
}

type configService struct{}

// NewConfigService creates a new ConfigService.
func NewConfigService() ConfigService {
	return &configService{}
}

func (s *configService) GetStatus(mcpServer *mcp.Server, proxyAddr string) (map[string]any, error) {
	mcpSessions := 0
	if mcpServer != nil {
		mcpSessions = mcpServer.ActiveSessions()
	}
	return map[string]any{
		"version":      config.Version,
		"proxy_addr":   proxyAddr,
		"mcp_sessions": mcpSessions,
		"mcp_enabled":  mcpServer != nil,
	}, nil
}

func (s *configService) GetConfig() (*model.Config, error) {
	return config.Get(), nil
}

func (s *configService) SaveConfig(cfg *model.Config) error {
	return config.Save(cfg)
}
