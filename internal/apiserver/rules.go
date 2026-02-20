package apiserver

import (
	"glance/internal/model"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleCreateRule(c *fiber.Ctx) error {
	rule := new(model.Rule)
	if err := c.BodyParser(rule); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.services.Rule.Create(rule); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(rule)
}

func (s *Server) handleListRules(c *fiber.Ctx) error {
	return c.JSON(s.services.Rule.GetAll())
}

func (s *Server) handleUpdateRule(c *fiber.Ctx) error {
	id := c.Params("id")
	rule := new(model.Rule)
	if err := c.BodyParser(rule); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.services.Rule.Update(id, rule); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(rule)
}

func (s *Server) handleDeleteRule(c *fiber.Ctx) error {
	id := c.Params("id")
	s.services.Rule.Delete(id)
	return c.SendStatus(fiber.StatusNoContent)
}
