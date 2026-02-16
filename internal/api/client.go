package api

import (
	"agent-proxy/internal/client"
	"github.com/gofiber/fiber/v2"
)

func (s *APIServer) handleLaunchChromium(c *fiber.Ctx) error {
	err := client.LaunchChromium(s.proxyAddr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "launched"})
}

func (s *APIServer) handleListJavaProcesses(c *fiber.Ctx) error {
	procs, err := client.ListJavaProcesses()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(procs)
}

func (s *APIServer) registerClientRoutes() {
	s.app.Post("/api/client/chromium", s.handleLaunchChromium)
	s.app.Get("/api/client/java/processes", s.handleListJavaProcesses)
}
