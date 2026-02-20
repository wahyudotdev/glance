package apiserver

import (
	"glance/internal/model"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleListScenarios(c *fiber.Ctx) error {
	scenarios, err := s.services.Scenario.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(scenarios)
}

func (s *Server) handleGetScenario(c *fiber.Ctx) error {
	id := c.Params("id")
	scenario, err := s.services.Scenario.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Scenario not found"})
	}
	return c.JSON(scenario)
}

func (s *Server) handleCreateScenario(c *fiber.Ctx) error {
	scenario := new(model.Scenario)
	if err := c.BodyParser(scenario); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.services.Scenario.Create(scenario); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(scenario)
}

func (s *Server) handleUpdateScenario(c *fiber.Ctx) error {
	id := c.Params("id")
	scenario := new(model.Scenario)
	if err := c.BodyParser(scenario); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.services.Scenario.Update(id, scenario); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(scenario)
}

func (s *Server) handleDeleteScenario(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := s.services.Scenario.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
