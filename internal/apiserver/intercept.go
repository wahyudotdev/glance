package apiserver

import (
	"encoding/json"
	"glance/internal/proxy"
	"glance/internal/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleContinueRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	var data struct {
		Method  string              `json:"method"`
		URL     string              `json:"url"`
		Headers map[string][]string `json:"headers"`
		Body    string              `json:"body"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	h := http.Header{}
	for k, vs := range data.Headers {
		for _, v := range vs {
			h.Add(k, v)
		}
	}

	params := service.ContinueRequestParams{
		Method:  data.Method,
		URL:     data.URL,
		Headers: h,
		Body:    data.Body,
	}

	if err := s.services.Intercept.ContinueRequest(id, params); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "resumed"})
}

func (s *Server) handleContinueResponse(c *fiber.Ctx) error {
	id := c.Params("id")
	var data struct {
		Status  int                 `json:"status"`
		Headers map[string][]string `json:"headers"`
		Body    string              `json:"body"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	h := http.Header{}
	for k, vs := range data.Headers {
		for _, v := range vs {
			h.Add(k, v)
		}
	}

	params := service.ContinueResponseParams{
		Status:  data.Status,
		Headers: h,
		Body:    data.Body,
	}

	if err := s.services.Intercept.ContinueResponse(id, params); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "resumed"})
}

func (s *Server) handleAbortRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := s.services.Intercept.Abort(id); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "aborted"})
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
	s.Hub.BroadcastData(data)
}
