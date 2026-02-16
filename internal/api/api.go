package api

import (
	"agent-proxy/internal/config"
	"agent-proxy/internal/interceptor"
	"agent-proxy/internal/model"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

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
	s.app.Post("/api/request/execute", s.handleExecuteRequest)
	s.app.Post("/api/request/execute", s.handleExecuteRequest)

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
	cfg := config.Get()
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", cfg.DefaultPageSize)

	offset := (page - 1) * pageSize
	entries, total := s.store.GetPage(offset, pageSize)

	return c.JSON(fiber.Map{
		"entries":  entries,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (s *APIServer) handleClearTraffic(c *fiber.Ctx) error {
	s.store.ClearEntries()
	return c.SendStatus(fiber.StatusNoContent)
}

func (s *APIServer) handleGetConfig(c *fiber.Ctx) error {
	return c.JSON(config.Get())
}

func (s *APIServer) handleSaveConfig(c *fiber.Ctx) error {
	cfg := new(model.Config)
	if err := c.BodyParser(cfg); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := config.Save(cfg); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cfg)
}

func (s *APIServer) handleExecuteRequest(c *fiber.Ctx) error {
	type ExecuteRequest struct {
		Method  string              `json:"method"`
		URL     string              `json:"url"`
		Headers map[string][]string `json:"headers"`
		Body    string              `json:"body"`
	}

	reqData := new(ExecuteRequest)
	if err := c.BodyParser(reqData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Prepare the HTTP request
	req, err := http.NewRequest(reqData.Method, reqData.URL, strings.NewReader(reqData.Body))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Copy headers
	for k, vs := range reqData.Headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	// Capture the request start
	entry, _ := interceptor.NewEntry(req)

	// Execute
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	start := time.Now()
	resp, err := client.Do(req)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer resp.Body.Close()

	// Capture the response
	entry.Duration = time.Since(start)
	entry.Status = resp.StatusCode
	entry.ResponseHeaders = resp.Header.Clone()
	bodyBytes, _ := io.ReadAll(resp.Body)
	entry.ResponseBody = string(bodyBytes)

	// Detect and encode image if necessary (reuse logic from interceptor)
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "image/") {
		encoded := base64.StdEncoding.EncodeToString(bodyBytes)
		entry.ResponseBody = fmt.Sprintf("data:%s;base64,%s", contentType, encoded)
	}

	// Save to store
	s.store.AddEntry(entry)

	// Also broadcast via WebSocket
	s.Hub.Broadcast(entry)

	return c.JSON(entry)
}

func (s *APIServer) Listen(addr string) error {
	return s.app.Listen(addr)
}
