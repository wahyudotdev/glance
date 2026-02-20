package apiserver

import (
	"glance/internal/service"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleExecuteRequest(c *fiber.Ctx) error {
	var params service.ExecuteRequestParams
	if err := c.BodyParser(&params); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	entry, err := s.services.Request.Execute(params)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Also broadcast via WebSocket
	s.Hub.Broadcast(entry)

	return c.JSON(entry)
}
