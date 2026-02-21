package apiserver

import (
	"glance/internal/model"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type mockClientService struct{}

func (m *mockClientService) LaunchChromium(_ string) error                      { return nil }
func (m *mockClientService) ListJavaProcesses() ([]model.JavaProcess, error)    { return nil, nil }
func (m *mockClientService) InterceptJava(_, _ string) error                    { return nil }
func (m *mockClientService) GetTerminalSetupScript(_ string) string             { return "script" }
func (m *mockClientService) ListAndroidDevices() ([]model.AndroidDevice, error) { return nil, nil }
func (m *mockClientService) InterceptAndroid(_, _ string) error                 { return nil }
func (m *mockClientService) ClearAndroid(_, _ string) error                     { return nil }
func (m *mockClientService) PushAndroidCert(_ string) error                     { return nil }

func TestHandleLaunchChromium(t *testing.T) {
	app := fiber.New()
	svc := &mockClientService{}
	s := &Server{
		services: Services{Client: svc},
		app:      app,
	}
	app.Post("/api/client/chromium", s.handleLaunchChromium)

	req := httptest.NewRequest("POST", "/api/client/chromium", nil)
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleListJavaProcesses(t *testing.T) {
	app := fiber.New()
	svc := &mockClientService{}
	s := &Server{
		services: Services{Client: svc},
		app:      app,
	}
	app.Get("/api/client/java/processes", s.handleListJavaProcesses)

	req := httptest.NewRequest("GET", "/api/client/java/processes", nil)
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleInterceptJava(t *testing.T) {
	app := fiber.New()
	svc := &mockClientService{}
	s := &Server{
		services: Services{Client: svc},
		app:      app,
	}
	app.Post("/api/client/java/intercept/:pid", s.handleInterceptJava)

	req := httptest.NewRequest("POST", "/api/client/java/intercept/123", nil)
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleTerminalSetup(t *testing.T) {
	app := fiber.New()
	svc := &mockClientService{}
	s := &Server{
		services: Services{Client: svc},
		app:      app,
	}
	app.Get("/api/client/terminal/setup", s.handleTerminalSetup)

	req := httptest.NewRequest("GET", "/api/client/terminal/setup", nil)
	resp, _ := app.Test(req)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleAndroidHandlers(t *testing.T) {
	app := fiber.New()
	svc := &mockClientService{}
	s := &Server{
		services: Services{Client: svc},
		app:      app,
	}
	app.Get("/api/client/android/devices", s.handleListAndroidDevices)
	app.Post("/api/client/android/intercept/:id", s.handleInterceptAndroid)
	app.Post("/api/client/android/clear/:id", s.handleClearAndroid)
	app.Post("/api/client/android/push-cert/:id", s.handlePushAndroidCert)

	t.Run("List", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/client/android/devices", nil)
		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Intercept", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/client/android/intercept/dev1", nil)
		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Clear", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/client/android/clear/dev1", nil)
		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PushCert", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/client/android/push-cert/dev1", nil)
		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}
