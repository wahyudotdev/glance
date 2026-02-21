package apiserver

import (
	"bytes"
	"glance/internal/model"
	"glance/internal/service"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type mockRequestService struct{}

func (m *mockRequestService) Execute(_ service.ExecuteRequestParams) (*model.TrafficEntry, error) {
	return &model.TrafficEntry{ID: "123", Status: 200}, nil
}

func TestHandleExecuteRequest(t *testing.T) {
	app := fiber.New()
	svc := &mockRequestService{}
	// Also need a hub for broadcasting
	hub := NewHub()
	go hub.Run()

	s := &Server{
		services: Services{Request: svc},
		app:      app,
		Hub:      hub,
	}
	app.Post("/api/request/execute", s.handleExecuteRequest)

	body := `{"method":"GET", "url":"http://test.com"}`
	req := httptest.NewRequest("POST", "/api/request/execute", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
