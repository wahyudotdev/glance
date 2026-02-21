package apiserver

import (
	"bytes"
	"glance/internal/model"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHandleListRules(t *testing.T) {
	app := fiber.New()
	svc := &mockRuleService{rules: []*model.Rule{{ID: "1"}}}
	s := &Server{
		services: Services{Rule: svc},
		app:      app,
	}
	app.Get("/api/rules", s.handleListRules)

	req := httptest.NewRequest("GET", "/api/rules", nil)
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleCreateRule(t *testing.T) {
	app := fiber.New()
	svc := &mockRuleService{}
	s := &Server{
		services: Services{Rule: svc},
		app:      app,
	}
	app.Post("/api/rules", s.handleCreateRule)

	body := `{"type":"mock", "url_pattern":"/api"}`
	req := httptest.NewRequest("POST", "/api/rules", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleUpdateRule(t *testing.T) {
	app := fiber.New()
	svc := &mockRuleService{}
	s := &Server{
		services: Services{Rule: svc},
		app:      app,
	}
	app.Put("/api/rules/:id", s.handleUpdateRule)

	body := `{"type":"mock", "url_pattern":"/new"}`
	req := httptest.NewRequest("PUT", "/api/rules/123", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleDeleteRule(t *testing.T) {
	app := fiber.New()
	svc := &mockRuleService{}
	s := &Server{
		services: Services{Rule: svc},
		app:      app,
	}
	app.Delete("/api/rules/:id", s.handleDeleteRule)

	req := httptest.NewRequest("DELETE", "/api/rules/123", nil)
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 204 {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}
