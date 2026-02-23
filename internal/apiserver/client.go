package apiserver

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleLaunchChromium(c *fiber.Ctx) error {
	err := s.services.Client.LaunchChromium(s.proxyAddr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "launched"})
}

func (s *Server) handleListJavaProcesses(c *fiber.Ctx) error {
	procs, err := s.services.Client.ListJavaProcesses()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(procs)
}

func (s *Server) handleInterceptJava(c *fiber.Ctx) error {
	pid := c.Params("pid")
	err := s.services.Client.InterceptJava(pid, s.proxyAddr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "intercepted"})
}

func (s *Server) handleTerminalSetup(c *fiber.Ctx) error {
	script := s.services.Client.GetTerminalSetupScript(s.proxyAddr)
	return c.SendString(script)
}

func (s *Server) handleListAndroidDevices(c *fiber.Ctx) error {
	devices, err := s.services.Client.ListAndroidDevices()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(devices)
}

func (s *Server) handleInterceptAndroid(c *fiber.Ctx) error {
	deviceID := c.Params("id")
	err := s.services.Client.InterceptAndroid(deviceID, s.proxyAddr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "intercepted"})
}

func (s *Server) handleClearAndroid(c *fiber.Ctx) error {
	deviceID := c.Params("id")
	err := s.services.Client.ClearAndroid(deviceID, s.proxyAddr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "cleared"})
}

func (s *Server) handlePushAndroidCert(c *fiber.Ctx) error {
	deviceID := c.Params("id")
	err := s.services.Client.PushAndroidCert(deviceID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "pushed"})
}

func (s *Server) handleListDockerContainers(c *fiber.Ctx) error {
	containers, err := s.services.Client.ListDockerContainers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(containers)
}

func (s *Server) handleInterceptDocker(c *fiber.Ctx) error {
	containerID := c.Params("id")
	err := s.services.Client.InterceptDocker(containerID, s.proxyAddr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "intercepting"})
}

func (s *Server) handleStopInterceptDocker(c *fiber.Ctx) error {
	containerID := c.Params("id")
	err := s.services.Client.StopInterceptDocker(containerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "stopped"})
}

func (s *Server) registerClientRoutes() {
	s.app.Post("/api/client/chromium", s.handleLaunchChromium)
	s.app.Get("/api/client/java/processes", s.handleListJavaProcesses)
	s.app.Post("/api/client/java/intercept/:pid", s.handleInterceptJava)

	s.app.Get("/api/client/android/devices", s.handleListAndroidDevices)
	s.app.Post("/api/client/android/intercept/:id", s.handleInterceptAndroid)
	s.app.Post("/api/client/android/clear/:id", s.handleClearAndroid)
	s.app.Post("/api/client/android/push-cert/:id", s.handlePushAndroidCert)

	s.app.Get("/api/client/docker/containers", s.handleListDockerContainers)
	s.app.Post("/api/client/docker/intercept/:id", s.handleInterceptDocker)
	s.app.Post("/api/client/docker/stop/:id", s.handleStopInterceptDocker)

	s.app.Get("/api/client/terminal/setup", s.handleTerminalSetup)
	s.app.Get("/setup", s.handleTerminalSetup)
}
