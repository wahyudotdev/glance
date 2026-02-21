package apiserver

import (
	"bytes"
	"encoding/json"
	"glance/internal/model"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHandleStatus(t *testing.T) {
	app := fiber.New()
	svc := &mockConfigService{status: map[string]any{"version": "1.0.0"}}
	s := &Server{
		services: Services{Config: svc},
		app:      app,
	}
	app.Get("/api/status", s.handleStatus)

	req := httptest.NewRequest("GET", "/api/status", nil)
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var body map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body["version"] != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %v", body["version"])
	}
}

func TestHandleGetConfig(t *testing.T) {
	app := fiber.New()
	svc := &mockConfigService{cfg: &model.Config{ProxyAddr: ":8000"}}
	s := &Server{
		services: Services{Config: svc},
		app:      app,
	}
	app.Get("/api/config", s.handleGetConfig)

	req := httptest.NewRequest("GET", "/api/config", nil)
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleSaveConfig(t *testing.T) {
	app := fiber.New()
	svc := &mockConfigService{}
	s := &Server{
		services: Services{Config: svc},
		app:      app,
	}
	app.Post("/api/config", s.handleSaveConfig)

	body := `{"proxy_addr":":9000"}`
	req := httptest.NewRequest("POST", "/api/config", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
