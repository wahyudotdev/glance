package apiserver

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type mockCAService struct{}

func (m *mockCAService) GetCACert() []byte { return []byte("cert") }

func TestHandleDownloadCA(t *testing.T) {
	app := fiber.New()
	svc := &mockCAService{}
	s := &Server{
		services: Services{CA: svc},
		app:      app,
	}
	app.Get("/api/ca/cert", s.handleDownloadCA)

	req := httptest.NewRequest("GET", "/api/ca/cert", nil)
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
