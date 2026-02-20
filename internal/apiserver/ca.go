package apiserver

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleDownloadCA(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-x509-ca-cert")
	c.Set("Content-Disposition", `attachment; filename="glance-ca.crt"`)
	return c.Send(s.services.CA.GetCACert())
}

func (s *Server) registerCARoutes() {
	s.app.Get("/api/ca/cert", s.handleDownloadCA)
}
