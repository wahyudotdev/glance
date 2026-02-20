package apiserver

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleTraffic(c *fiber.Ctx) error {
	cfg, _ := s.services.Config.GetConfig()
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", cfg.DefaultPageSize)

	offset := (page - 1) * pageSize
	entries, total := s.services.Traffic.GetPage(offset, pageSize)

	return c.JSON(fiber.Map{
		"entries":  entries,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (s *Server) handleClearTraffic(c *fiber.Ctx) error {
	s.services.Traffic.Clear()
	return c.SendStatus(fiber.StatusNoContent)
}
