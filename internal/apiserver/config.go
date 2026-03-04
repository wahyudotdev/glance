package apiserver

import (
	"glance/internal/model"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleStatus(c *fiber.Ctx) error {
	mcpSessions := 0
	mcpEnabled := false
	if s.mcp != nil {
		mcpSessions = s.mcp.ActiveSessions()
		mcpEnabled = true
	}

	status, err := s.services.Config.GetStatus(mcpSessions, mcpEnabled, s.proxyAddr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(status)
}

func (s *Server) handleGetConfig(c *fiber.Ctx) error {
	cfg, err := s.services.Config.GetConfig()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cfg)
}

func (s *Server) handleSaveConfig(c *fiber.Ctx) error {
	cfg := new(model.Config)
	if err := c.BodyParser(cfg); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := s.services.Config.SaveConfig(cfg); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cfg)
}
