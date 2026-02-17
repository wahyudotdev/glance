// Package apiserver implements the REST and WebSocket API for the dashboard and external clients.
package apiserver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"glance/internal/config"
	"glance/internal/interceptor"
	"glance/internal/mcp"
	"glance/internal/model"
	"glance/internal/proxy"
	"glance/internal/repository"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

// Server manages the HTTP and WebSocket endpoints for the application.
type Server struct {
	store        *interceptor.TrafficStore
	proxy        *proxy.Proxy
	mcp          *mcp.Server
	scenarioRepo repository.ScenarioRepository
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

func (s *Server) handleStatus(c *fiber.Ctx) error {
	mcpSessions := 0
	if s.mcp != nil {
		mcpSessions = s.mcp.ActiveSessions()
	}
	return c.JSON(fiber.Map{
		"proxy_addr":   s.proxyAddr,
		"mcp_sessions": mcpSessions,
		"mcp_enabled":  s.mcp != nil,
	})
}

func (s *Server) handleTraffic(c *fiber.Ctx) error {
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

func (s *Server) handleClearTraffic(c *fiber.Ctx) error {
	s.store.ClearEntries()
	return c.SendStatus(fiber.StatusNoContent)
}

func (s *Server) handleGetConfig(c *fiber.Ctx) error {
	return c.JSON(config.Get())
}

func (s *Server) handleSaveConfig(c *fiber.Ctx) error {
	cfg := new(model.Config)
	if err := c.BodyParser(cfg); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := config.Save(cfg); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cfg)
}

func (s *Server) handleExecuteRequest(c *fiber.Ctx) error {
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
	entry.ModifiedBy = "editor"

	// Execute
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	start := time.Now()
	resp, err := client.Do(req)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer resp.Body.Close() //nolint:errcheck

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

func (s *Server) handleContinueRequest(c *fiber.Ctx) error {
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

func (s *Server) handleContinueResponse(c *fiber.Ctx) error {
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

func (s *Server) handleAbortRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	success := s.proxy.AbortRequest(id)
	if !success {
		return c.Status(404).JSON(fiber.Map{"error": "Intercepted request not found or already released"})
	}
	return c.JSON(fiber.Map{"status": "aborted"})
}

func (s *Server) handleCreateRule(c *fiber.Ctx) error {
	rule := new(model.Rule)
	if err := c.BodyParser(rule); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if rule.ID == "" {
		rule.ID = uuid.New().String()
	}

	s.proxy.Engine.AddRule(rule)
	return c.JSON(rule)
}

func (s *Server) handleListRules(c *fiber.Ctx) error {
	return c.JSON(s.proxy.Engine.GetRules())
}

func (s *Server) handleUpdateRule(c *fiber.Ctx) error {
	id := c.Params("id")
	rule := new(model.Rule)
	if err := c.BodyParser(rule); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	rule.ID = id
	s.proxy.Engine.UpdateRule(rule)
	return c.JSON(rule)
}

func (s *Server) handleDeleteRule(c *fiber.Ctx) error {
	id := c.Params("id")
	s.proxy.Engine.DeleteRule(id)
	return c.SendStatus(fiber.StatusNoContent)
}

// BroadcastIntercept sends an interception event to all connected WebSocket clients.
func (s *Server) BroadcastIntercept(bp *proxy.Breakpoint) {
	msg := fiber.Map{
		"type":           "intercepted",
		"intercept_type": bp.Type,
		"id":             bp.ID,
		"entry":          bp.Entry,
	}

	data, _ := json.Marshal(msg)
	s.Hub.mu.Lock()
	for client := range s.Hub.clients {
		_ = client.WriteMessage(websocket.TextMessage, data) //nolint:errcheck
	}
	s.Hub.mu.Unlock()
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

	// Create a standard http.ServeMux to handle MCP correctly (real SSE flushing)
	mux := http.NewServeMux()

	if s.mcp != nil {
		// Use the official SDK's StreamableHTTPHandler (handles GET and POST)
		mcpHandler := s.mcp.GetStreamableHandler()

		// CORS and Logging Middleware for standard net/http
		mcpWrapper := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Handle CORS
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-MCP-Protocol-Version")

				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusNoContent)
					return
				}

				// Log MCP Request
				log.Printf("\033[35m[MCP]\033[0m %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

				h.ServeHTTP(w, r)
			})
		}

		// Mount official handler with wrapper
		mux.Handle("/mcp", mcpWrapper(mcpHandler))
		fmt.Printf("\033[32m[✓]\033[0m MCP server (Streamable HTTP) unified on \033[34m\033[1mhttp://%s/mcp\033[0m\n", displayAddr)
	}

	// Use adaptor to mount the entire Fiber app on the remaining routes
	mux.HandleFunc("/", adaptor.FiberApp(s.app))

	// Use a standard http server to host the mux
	srv := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return actualAddr, srv.Serve(ln)
}
