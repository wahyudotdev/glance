package apiserver

import (
	"glance/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *Server) handleListScenarios(c *fiber.Ctx) error {
	scenarios, err := s.scenarioRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(scenarios)
}

func (s *Server) handleGetScenario(c *fiber.Ctx) error {
	id := c.Params("id")
	scenario, err := s.scenarioRepo.GetByID(id)
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

	if scenario.ID == "" {
		scenario.ID = uuid.New().String()
	}
	scenario.CreatedAt = time.Now()

	if err := s.scenarioRepo.Add(scenario); err != nil {
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

	scenario.ID = id
	if err := s.scenarioRepo.Update(scenario); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(scenario)
}

func (s *Server) handleDeleteScenario(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := s.scenarioRepo.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
