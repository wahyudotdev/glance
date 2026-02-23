package apiserver

import (
	"errors"
	"glance/internal/model"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type mockClientService struct {
	err error
}

func (m *mockClientService) LaunchChromium(_ string) error { return m.err }
func (m *mockClientService) ListJavaProcesses() ([]model.JavaProcess, error) {
	return nil, m.err
}
func (m *mockClientService) InterceptJava(_, _ string) error        { return m.err }
func (m *mockClientService) GetTerminalSetupScript(_ string) string { return "script" }
func (m *mockClientService) ListAndroidDevices() ([]model.AndroidDevice, error) {
	return nil, m.err
}
func (m *mockClientService) InterceptAndroid(_, _ string) error { return m.err }
func (m *mockClientService) ClearAndroid(_, _ string) error     { return m.err }
func (m *mockClientService) PushAndroidCert(_ string) error     { return m.err }
func (m *mockClientService) ListDockerContainers() ([]model.DockerContainer, error) {
	return nil, m.err
}
func (m *mockClientService) InterceptDocker(_, _ string) error  { return m.err }
func (m *mockClientService) StopInterceptDocker(_ string) error { return m.err }

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

	// Error case
	svc.err = errors.New("launch error")
	req = httptest.NewRequest("POST", "/api/client/chromium", nil)
	resp, _ = app.Test(req)
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != 500 {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
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

	// Error case
	svc.err = errors.New("list error")
	req = httptest.NewRequest("GET", "/api/client/java/processes", nil)
	resp, _ = app.Test(req)
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != 500 {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
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

	// Error case
	svc.err = errors.New("intercept error")
	req = httptest.NewRequest("POST", "/api/client/java/intercept/123", nil)
	resp, _ = app.Test(req)
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != 500 {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
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
		svc.err = nil
		req := httptest.NewRequest("GET", "/api/client/android/devices", nil)
		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		svc.err = errors.New("error")
		resp, _ = app.Test(req)
		if resp.StatusCode != 500 {
			t.Errorf("Expected 500, got %d", resp.StatusCode)
		}
	})

	t.Run("Intercept", func(t *testing.T) {
		svc.err = nil
		req := httptest.NewRequest("POST", "/api/client/android/intercept/dev1", nil)
		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		svc.err = errors.New("error")
		resp, _ = app.Test(req)
		if resp.StatusCode != 500 {
			t.Errorf("Expected 500, got %d", resp.StatusCode)
		}
	})

	t.Run("Clear", func(t *testing.T) {
		svc.err = nil
		req := httptest.NewRequest("POST", "/api/client/android/clear/dev1", nil)
		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		svc.err = errors.New("error")
		resp, _ = app.Test(req)
		if resp.StatusCode != 500 {
			t.Errorf("Expected 500, got %d", resp.StatusCode)
		}
	})

	t.Run("PushCert", func(t *testing.T) {
		svc.err = nil
		req := httptest.NewRequest("POST", "/api/client/android/push-cert/dev1", nil)
		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		svc.err = errors.New("error")
		resp, _ = app.Test(req)
		if resp.StatusCode != 500 {
			t.Errorf("Expected 500, got %d", resp.StatusCode)
		}
	})
}

func TestHandleDockerHandlers(t *testing.T) {
	app := fiber.New()
	svc := &mockClientService{}
	s := &Server{
		services: Services{Client: svc},
		app:      app,
	}
	app.Get("/api/client/docker/containers", s.handleListDockerContainers)
	app.Post("/api/client/docker/intercept/:id", s.handleInterceptDocker)
	app.Post("/api/client/docker/stop/:id", s.handleStopInterceptDocker)

	t.Run("List", func(t *testing.T) {
		svc.err = nil
		req := httptest.NewRequest("GET", "/api/client/docker/containers", nil)
		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		svc.err = errors.New("error")
		resp, _ = app.Test(req)
		if resp.StatusCode != 500 {
			t.Errorf("Expected 500, got %d", resp.StatusCode)
		}
	})

	t.Run("Intercept", func(t *testing.T) {
		svc.err = nil
		req := httptest.NewRequest("POST", "/api/client/docker/intercept/cont1", nil)
		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		svc.err = errors.New("error")
		resp, _ = app.Test(req)
		if resp.StatusCode != 500 {
			t.Errorf("Expected 500, got %d", resp.StatusCode)
		}
	})

	t.Run("Stop", func(t *testing.T) {
		svc.err = nil
		req := httptest.NewRequest("POST", "/api/client/docker/stop/cont1", nil)
		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		svc.err = errors.New("error")
		resp, _ = app.Test(req)
		if resp.StatusCode != 500 {
			t.Errorf("Expected 500, got %d", resp.StatusCode)
		}
	})
}
