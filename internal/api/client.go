package api

import (
	"agent-proxy/internal/client"
	"github.com/gofiber/fiber/v2"
)

func (s *APIServer) handleLaunchChromium(c *fiber.Ctx) error {
	// Assuming proxy is on :8080 as per requirements
	err := client.LaunchChromium("127.0.0.1:8080")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "launched"})
}

func (s *APIServer) registerClientRoutes() {
	s.app.Post("/api/client/chromium", s.handleLaunchChromium)
}
