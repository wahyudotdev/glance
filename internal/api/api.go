package api

import (
	"agent-proxy/internal/config"
	"agent-proxy/internal/interceptor"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

type APIServer struct {
	store       *interceptor.TrafficStore
	app         *fiber.App
	proxyAddr   string
	restartChan chan bool
	Hub         *Hub
}

func NewAPIServer(store *interceptor.TrafficStore, proxyAddr string) *APIServer {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	hub := NewHub()
	go hub.Run()

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "*",
	}))

	// Add request logging for API debugging
	app.Use(func(c *fiber.Ctx) error {
		if strings.HasPrefix(c.Path(), "/api") {
			log.Printf("[API] %s %s", c.Method(), c.Path())
		}
		return c.Next()
	})

	return &APIServer{
		store:       store,
		app:         app,
		proxyAddr:   proxyAddr,
		restartChan: make(chan bool, 1),
		Hub:         hub,
	}
}

func (s *APIServer) RegisterRoutes() {
	s.app.Get("/api/status", s.handleStatus)
	s.app.Get("/api/traffic", s.handleTraffic)
	s.app.Delete("/api/traffic", s.handleClearTraffic)
	s.app.Get("/api/config", s.handleGetConfig)
	s.app.Post("/api/config", s.handleSaveConfig)

	// WebSocket for real-time traffic
	s.app.Get("/ws/traffic", websocket.New(func(c *websocket.Conn) {
		s.Hub.register <- c
		defer func() {
			s.Hub.unregister <- c
		}()

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}
		}
	}))

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
