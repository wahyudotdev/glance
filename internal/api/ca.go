package api

import (
	"github.com/elazarl/goproxy"
	"github.com/gofiber/fiber/v2"
)

func (s *APIServer) handleDownloadCA(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-x509-ca-cert")
	c.Set("Content-Disposition", `attachment; filename="agent-proxy-ca.crt"`)
	return c.Send(goproxy.CA_CERT)
}

func (s *APIServer) registerCARoutes() {
	s.app.Get("/api/ca/cert", s.handleDownloadCA)
}
