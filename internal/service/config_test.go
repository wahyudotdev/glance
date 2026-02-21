package service

import (
	"glance/internal/config"
	"glance/internal/mcp"
	"glance/internal/model"
	"testing"
)

type mockConfigRepo struct {
	cfg *model.Config
}

func (m *mockConfigRepo) Get() (*model.Config, error) {
	return m.cfg, nil
}

func (m *mockConfigRepo) Save(c *model.Config) error {
	m.cfg = c
	return nil
}

func TestConfigService_Status_MCP(t *testing.T) {
	config.Init(&mockConfigRepo{cfg: &model.Config{}})
	svc := NewConfigService()

	// Minimal MCP server setup
	ms := mcp.NewServer(nil, nil, ":0", nil)
	status, _ := svc.GetStatus(ms, ":8000")
	if status["mcp_enabled"] != true {
		t.Error("Expected mcp_enabled true")
	}
}

func TestConfigService(t *testing.T) {
	repo := &mockConfigRepo{cfg: &model.Config{ProxyAddr: ":8000"}}
	config.Init(repo)
	svc := NewConfigService()

	// Test GetStatus
	status, err := svc.GetStatus(nil, ":8000")
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status["proxy_addr"] != ":8000" {
		t.Errorf("Expected proxy_addr :8000, got %v", status["proxy_addr"])
	}

	// Test GetConfig
	cfg, err := svc.GetConfig()
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}
	if cfg.ProxyAddr != ":8000" {
		t.Errorf("Expected ProxyAddr :8000, got %s", cfg.ProxyAddr)
	}

	// Test SaveConfig
	cfg.ProxyAddr = ":9000"
	err = svc.SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}
	if repo.cfg.ProxyAddr != ":9000" {
		t.Errorf("Expected saved ProxyAddr :9000, got %s", repo.cfg.ProxyAddr)
	}
}
