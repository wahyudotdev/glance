package apiserver

import (
	"github.com/elazarl/goproxy"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleDownloadCA(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-x509-ca-cert")
	c.Set("Content-Disposition", `attachment; filename="agent-proxy-ca.crt"`)
	return c.Send(goproxy.CA_CERT)
}

func (s *Server) registerCARoutes() {
	s.app.Get("/api/ca/cert", s.handleDownloadCA)
}
