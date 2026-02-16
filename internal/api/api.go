package api

import (
	"agent-proxy/internal/config"
	"agent-proxy/internal/interceptor"
	"agent-proxy/internal/model"
	"agent-proxy/internal/proxy"
	"agent-proxy/internal/rules"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type APIServer struct {
	store       *interceptor.TrafficStore
	proxy       *proxy.Proxy
	app         *fiber.App
	proxyAddr   string
	restartChan chan bool
	Hub         *Hub
}

func NewAPIServer(store *interceptor.TrafficStore, p *proxy.Proxy, proxyAddr string) *APIServer {
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
		proxy:       p,
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
	s.app.Get("/api/rules", s.handleListRules)
	s.app.Post("/api/rules", s.handleCreateRule)
	s.app.Put("/api/rules/:id", s.handleUpdateRule)
	s.app.Delete("/api/rules/:id", s.handleDeleteRule)
	s.app.Post("/api/intercept/continue/:id", s.handleContinueRequest)
	s.app.Post("/api/intercept/response/continue/:id", s.handleContinueResponse)
	s.app.Post("/api/intercept/abort/:id", s.handleAbortRequest)

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

func (s *APIServer) handleContinueRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	type ContinueRequest struct {
		Method  string              `json:"method"`
		URL     string              `json:"url"`
		Headers map[string][]string `json:"headers"`
		Body    string              `json:"body"`
	}

	reqData := new(ContinueRequest)
	if err := c.BodyParser(reqData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	h := http.Header{}
	for k, vs := range reqData.Headers {
		for _, v := range vs {
			h.Add(k, v)
		}
	}

	success := s.proxy.ContinueRequest(id, reqData.Method, reqData.URL, h, reqData.Body)
	if !success {
		return c.Status(404).JSON(fiber.Map{"error": "Intercepted request not found or already released"})
	}

	return c.JSON(fiber.Map{"status": "resumed"})
}

func (s *APIServer) handleContinueResponse(c *fiber.Ctx) error {
	id := c.Params("id")
	type ContinueResponse struct {
		Status  int                 `json:"status"`
		Headers map[string][]string `json:"headers"`
		Body    string              `json:"body"`
	}

	resData := new(ContinueResponse)
	if err := c.BodyParser(resData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	h := http.Header{}
	for k, vs := range resData.Headers {
		for _, v := range vs {
			h.Add(k, v)
		}
	}

	success := s.proxy.ContinueResponse(id, resData.Status, h, resData.Body)
	if !success {
		return c.Status(404).JSON(fiber.Map{"error": "Intercepted response not found or already released"})
	}

	return c.JSON(fiber.Map{"status": "resumed"})
}

func (s *APIServer) handleAbortRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	success := s.proxy.AbortRequest(id)
	if !success {
		return c.Status(404).JSON(fiber.Map{"error": "Intercepted request not found or already released"})
	}
	return c.JSON(fiber.Map{"status": "aborted"})
}

func (s *APIServer) handleCreateRule(c *fiber.Ctx) error {
	rule := new(rules.Rule)
	if err := c.BodyParser(rule); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if rule.ID == "" {
		rule.ID = uuid.New().String()
	}

	s.proxy.Engine.AddRule(rule)
	return c.JSON(rule)
}

func (s *APIServer) handleListRules(c *fiber.Ctx) error {
	return c.JSON(s.proxy.Engine.GetRules())
}

func (s *APIServer) handleUpdateRule(c *fiber.Ctx) error {
	id := c.Params("id")
	rule := new(rules.Rule)
	if err := c.BodyParser(rule); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	rule.ID = id
	s.proxy.Engine.UpdateRule(rule)
	return c.JSON(rule)
}

func (s *APIServer) handleDeleteRule(c *fiber.Ctx) error {
	id := c.Params("id")
	s.proxy.Engine.DeleteRule(id)
	return c.SendStatus(fiber.StatusNoContent)
}

func (s *APIServer) BroadcastIntercept(bp *proxy.Breakpoint) {

	msg := fiber.Map{

		"type": "intercepted",

		"intercept_type": bp.Type,

		"id": bp.ID,

		"entry": bp.Entry,
	}

	data, _ := json.Marshal(msg)
	s.Hub.mu.Lock()
	for client := range s.Hub.clients {
		client.WriteMessage(websocket.TextMessage, data)
	}
	s.Hub.mu.Unlock()
}

func (s *APIServer) Listen(addr string) error {
	return s.app.Listen(addr)
}
