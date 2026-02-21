package config

import (
	"fmt"
	"glance/internal/model"
	"testing"
)

type mockRepo struct {
	cfg *model.Config
}

func (m *mockRepo) Get() (*model.Config, error) { return m.cfg, nil }
func (m *mockRepo) Save(c *model.Config) error  { m.cfg = c; return nil }

type errorRepo struct{}

func (e *errorRepo) Get() (*model.Config, error) { return nil, fmt.Errorf("error") }
func (e *errorRepo) Save(_ *model.Config) error  { return fmt.Errorf("error") }

func TestConfig_Fallback(t *testing.T) {
	Init(&errorRepo{})
	cfg := Get()
	if cfg.ProxyAddr != ":8000" {
		t.Errorf("Expected fallback proxy addr :8000, got %s", cfg.ProxyAddr)
	}
}

func TestConfig(t *testing.T) {
	repo := &mockRepo{cfg: &model.Config{ProxyAddr: ":8000"}}
	Init(repo)

	cfg := Get()
	if cfg.ProxyAddr != ":8000" {
		t.Errorf("Expected :8000, got %s", cfg.ProxyAddr)
	}

	cfg.ProxyAddr = ":9000"
	_ = Save(cfg)
	if repo.cfg.ProxyAddr != ":9000" {
		t.Errorf("Expected :9000, got %s", repo.cfg.ProxyAddr)
	}
}
