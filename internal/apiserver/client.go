package apiserver

import (
	"glance/internal/client"
	"net"

	"github.com/elazarl/goproxy"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleLaunchChromium(c *fiber.Ctx) error {
	err := client.LaunchChromium(s.proxyAddr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "launched"})
}

func (s *Server) handleListJavaProcesses(c *fiber.Ctx) error {
	procs, err := client.ListJavaProcesses()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(procs)
}

func (s *Server) handleInterceptJava(c *fiber.Ctx) error {
	pid := c.Params("pid")
	err := client.BuildAndAttachAgent(pid, s.proxyAddr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "intercepted"})
}

func (s *Server) handleTerminalSetup(c *fiber.Ctx) error {
	script := client.GetTerminalSetupScript(s.proxyAddr)
	return c.SendString(script)
}

func (s *Server) handleListAndroidDevices(c *fiber.Ctx) error {
	devices, err := client.ListAndroidDevices()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(devices)
}

func (s *Server) handleInterceptAndroid(c *fiber.Ctx) error {
	deviceID := c.Params("id")
	_, port, _ := net.SplitHostPort(s.proxyAddr)
	if port == "" {
		port = "8000" // Fallback
	}

	err := client.ConfigureAndroidProxy(deviceID, port)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "intercepted"})
}

func (s *Server) handleClearAndroid(c *fiber.Ctx) error {
	deviceID := c.Params("id")
	_, port, _ := net.SplitHostPort(s.proxyAddr)
	if port == "" {
		port = "8000"
	}

	err := client.ClearAndroidProxy(deviceID, port)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "cleared"})
}

func (s *Server) handlePushAndroidCert(c *fiber.Ctx) error {
	deviceID := c.Params("id")
	err := client.PushCertToDevice(deviceID, []byte(goproxy.CA_CERT))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "pushed"})
}

func (s *Server) registerClientRoutes() {
	s.app.Post("/api/client/chromium", s.handleLaunchChromium)
	s.app.Get("/api/client/java/processes", s.handleListJavaProcesses)
	s.app.Post("/api/client/java/intercept/:pid", s.handleInterceptJava)

	s.app.Get("/api/client/android/devices", s.handleListAndroidDevices)
	s.app.Post("/api/client/android/intercept/:id", s.handleInterceptAndroid)
	s.app.Post("/api/client/android/clear/:id", s.handleClearAndroid)
	s.app.Post("/api/client/android/push-cert/:id", s.handlePushAndroidCert)

	s.app.Get("/api/client/terminal/setup", s.handleTerminalSetup)
	s.app.Get("/setup", s.handleTerminalSetup)
}
