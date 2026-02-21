package apiserver

import (
	"glance/internal/model"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHandleTraffic(t *testing.T) {
	app := fiber.New()
	svc := &mockTrafficService{entries: []*model.TrafficEntry{{ID: "1"}}}
	cfgSvc := &mockConfigService{cfg: &model.Config{DefaultPageSize: 10}}
	s := &Server{
		services: Services{Traffic: svc, Config: cfgSvc},
		app:      app,
	}
	app.Get("/api/traffic", s.handleTraffic)

	req := httptest.NewRequest("GET", "/api/traffic", nil)
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleClearTraffic(t *testing.T) {
	app := fiber.New()
	svc := &mockTrafficService{}
	s := &Server{
		services: Services{Traffic: svc},
		app:      app,
	}
	app.Delete("/api/traffic", s.handleClearTraffic)

	req := httptest.NewRequest("DELETE", "/api/traffic", nil)
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 204 {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}
