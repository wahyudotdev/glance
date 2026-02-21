package apiserver

import (
	"bytes"
	"glance/internal/model"
	"glance/internal/proxy"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHandleContinueRequest(t *testing.T) {
	app := fiber.New()
	svc := &mockInterceptService{}
	s := &Server{
		services: Services{Intercept: svc},
		app:      app,
	}
	app.Post("/api/intercept/continue/:id", s.handleContinueRequest)

	body := `{"method":"GET", "url":"http://test.com"}`
	req := httptest.NewRequest("POST", "/api/intercept/continue/123", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleAbortRequest(t *testing.T) {
	app := fiber.New()
	svc := &mockInterceptService{}
	s := &Server{
		services: Services{Intercept: svc},
		app:      app,
	}
	app.Post("/api/intercept/abort/:id", s.handleAbortRequest)

	req := httptest.NewRequest("POST", "/api/intercept/abort/123", nil)
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleContinueResponse(t *testing.T) {
	app := fiber.New()
	svc := &mockInterceptService{}
	s := &Server{
		services: Services{Intercept: svc},
		app:      app,
	}
	app.Post("/api/intercept/response/continue/:id", s.handleContinueResponse)

	body := `{"status":200, "body":"ok"}`
	req := httptest.NewRequest("POST", "/api/intercept/response/continue/123", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestBroadcastIntercept(_ *testing.T) {
	hub := NewHub()
	go hub.Run()
	s := &Server{Hub: hub}

	bp := &proxy.Breakpoint{
		ID:    "123",
		Type:  "request",
		Entry: &model.TrafficEntry{ID: "123"},
	}

	s.BroadcastIntercept(bp)
}
