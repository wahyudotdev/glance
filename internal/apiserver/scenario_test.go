package apiserver

import (
	"bytes"
	"glance/internal/model"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHandleListScenarios(t *testing.T) {
	app := fiber.New()
	svc := &mockScenarioService{scenarios: []*model.Scenario{{ID: "1", Name: "Test"}}}
	s := &Server{
		services: Services{Scenario: svc},
		app:      app,
	}
	app.Get("/api/scenarios", s.handleListScenarios)

	req := httptest.NewRequest("GET", "/api/scenarios", nil)
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleGetScenario(t *testing.T) {
	app := fiber.New()
	svc := &mockScenarioService{scenarios: []*model.Scenario{{ID: "123", Name: "Test"}}}
	s := &Server{
		services: Services{Scenario: svc},
		app:      app,
	}
	app.Get("/api/scenarios/:id", s.handleGetScenario)

	req := httptest.NewRequest("GET", "/api/scenarios/123", nil)
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleDeleteScenario(t *testing.T) {
	app := fiber.New()
	svc := &mockScenarioService{}
	s := &Server{
		services: Services{Scenario: svc},
		app:      app,
	}
	app.Delete("/api/scenarios/:id", s.handleDeleteScenario)

	req := httptest.NewRequest("DELETE", "/api/scenarios/123", nil)
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 204 {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}

func TestHandleCreateScenario(t *testing.T) {
	app := fiber.New()
	svc := &mockScenarioService{}
	s := &Server{
		services: Services{Scenario: svc},
		app:      app,
	}
	app.Post("/api/scenarios", s.handleCreateScenario)

	body := `{"name":"New Scenario"}`
	req := httptest.NewRequest("POST", "/api/scenarios", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleUpdateScenario(t *testing.T) {
	app := fiber.New()
	svc := &mockScenarioService{}
	s := &Server{
		services: Services{Scenario: svc},
		app:      app,
	}
	app.Put("/api/scenarios/:id", s.handleUpdateScenario)

	body := `{"name":"Updated Scenario"}`
	req := httptest.NewRequest("PUT", "/api/scenarios/123", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
