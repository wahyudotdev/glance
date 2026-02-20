// Package apiserver implements the REST and WebSocket API for the dashboard and external clients.
package apiserver

import (
	"fmt"
	"glance/internal/interceptor"
	"glance/internal/mcp"
	"glance/internal/proxy"
	"glance/internal/repository"
	"glance/internal/service"
	"log"
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

// Services holds the business logic services used by the API server.
type Services struct {
	Config    service.ConfigService
	Traffic   service.TrafficService
	Rule      service.RuleService
	Intercept service.InterceptService
	Request   service.RequestService
	Scenario  service.ScenarioService
	Client    service.ClientService
	CA        service.CAService
}

// Server manages the HTTP and WebSocket endpoints for the application.
type Server struct {
	store        *interceptor.TrafficStore
	proxy        *proxy.Proxy
	mcp          *mcp.Server
	scenarioRepo repository.ScenarioRepository
	services     Services
	app          *fiber.App
	proxyAddr    string
	restartChan  chan bool
	Hub          *Hub
}

// NewServer creates and initializes a new Server instance.
func NewServer(store *interceptor.TrafficStore, p *proxy.Proxy, proxyAddr string, mcpServer *mcp.Server, scenarioRepo repository.ScenarioRepository) *Server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	hub := NewHub()
	go hub.Run()

	// Initialize Services
	services := Services{
		Config:    service.NewConfigService(),
		Traffic:   service.NewTrafficService(store),
		Rule:      service.NewRuleService(p.Engine),
		Intercept: service.NewInterceptService(p),
		Request:   service.NewRequestService(store),
		Scenario:  service.NewScenarioService(scenarioRepo),
		Client:    service.NewClientService(),
		CA:        service.NewCAService(),
	}

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

	return &Server{
		store:        store,
		proxy:        p,
		mcp:          mcpServer,
		scenarioRepo: scenarioRepo,
		services:     services,
		app:          app,
		proxyAddr:    proxyAddr,
		restartChan:  make(chan bool, 1),
		Hub:          hub,
	}
}

// RegisterRoutes sets up all the API and static asset routes.
func (s *Server) RegisterRoutes() {
	s.app.Get("/api/status", s.handleStatus)
	s.app.Get("/api/traffic", s.handleTraffic)
	s.app.Delete("/api/traffic", s.handleClearTraffic)
	s.app.Get("/api/config", s.handleGetConfig)
	s.app.Post("/api/config", s.handleSaveConfig)
	s.app.Post("/api/request/execute", s.handleExecuteRequest)
	s.app.Get("/api/rules", s.handleListRules)
	s.app.Post("/api/rules", s.handleCreateRule)
	s.app.Put("/api/rules/:id", s.handleUpdateRule)
	s.app.Delete("/api/rules/:id", s.handleDeleteRule)
	s.app.Post("/api/intercept/continue/:id", s.handleContinueRequest)
	s.app.Post("/api/intercept/response/continue/:id", s.handleContinueResponse)
	s.app.Post("/api/intercept/abort/:id", s.handleAbortRequest)

	// Scenario routes
	s.app.Get("/api/scenarios", s.handleListScenarios)
	s.app.Get("/api/scenarios/:id", s.handleGetScenario)
	s.app.Post("/api/scenarios", s.handleCreateScenario)
	s.app.Put("/api/scenarios/:id", s.handleUpdateScenario)
	s.app.Delete("/api/scenarios/:id", s.handleDeleteScenario)

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

	s.registerStaticRoutes()
}

// Listen starts the API server on the provided address and returns the actual address it bound to.
func (s *Server) Listen(addr string) (string, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("API Port %s is in use, falling back to a random port...", addr)
		ln, err = net.Listen("tcp", ":0") //nolint:gosec
		if err != nil {
			return "", err
		}
	}

	actualAddr := ln.Addr().String()
	displayAddr := actualAddr
	if strings.HasPrefix(actualAddr, "[::]") {
		displayAddr = "localhost" + strings.TrimPrefix(actualAddr, "[::]")
	}

	fmt.Printf("\033[32m[✓]\033[0m API server running on \033[1m%s\033[0m\n", displayAddr)
	fmt.Printf("\033[32m[✓]\033[0m Dashboard available at \033[34m\033[1mhttp://%s\033[0m\n", displayAddr)

	return actualAddr, s.app.Listener(ln)
}
