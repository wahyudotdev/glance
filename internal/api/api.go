package api

import (
	"agent-proxy/internal/config"
	"agent-proxy/internal/interceptor"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type APIServer struct {
	store       *interceptor.TrafficStore
	app         *fiber.App
	proxyAddr   string
	restartChan chan bool
}

func NewAPIServer(store *interceptor.TrafficStore, proxyAddr string) *APIServer {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,OPTIONS,DELETE",
		AllowHeaders: "Content-Type",
	}))

	return &APIServer{
		store:       store,
		app:         app,
		proxyAddr:   proxyAddr,
		restartChan: make(chan bool, 1),
	}
}

func (s *APIServer) RegisterRoutes() {
	s.app.Get("/api/status", s.handleStatus)
	s.app.Get("/api/traffic", s.handleTraffic)
	s.app.Delete("/api/traffic", s.handleClearTraffic)
	s.app.Get("/api/config", s.handleGetConfig)
	s.app.Post("/api/config", s.handleSaveConfig)

	// Client integrations
	s.registerClientRoutes()

	s.registerCARoutes()

	// Register static files (SPA)

	s.registerStaticRoutes()
}

func (s *APIServer) handleStatus(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"proxy_addr": s.proxyAddr,
	})
}

func (s *APIServer) handleTraffic(c *fiber.Ctx) error {
	entries := s.store.GetEntries()
	return c.JSON(entries)
}

func (s *APIServer) handleClearTraffic(c *fiber.Ctx) error {
	s.store.ClearEntries()
	return c.SendStatus(fiber.StatusNoContent)
}

func (s *APIServer) handleGetConfig(c *fiber.Ctx) error {
	return c.JSON(config.Get())
}

func (s *APIServer) handleSaveConfig(c *fiber.Ctx) error {
	cfg := new(config.Config)
	if err := c.BodyParser(cfg); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := config.Save(cfg); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cfg)
}

func (s *APIServer) Listen(addr string) error {
	return s.app.Listen(addr)
}
