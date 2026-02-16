// Package mcp implements the Model Context Protocol server for AI agent integration.
package mcp

import (
	"agent-proxy/internal/interceptor"
	"agent-proxy/internal/model"
	"agent-proxy/internal/rules"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Server manages the MCP connection and tool registrations.
type Server struct {
	store     *interceptor.TrafficStore
	engine    *rules.Engine
	proxyAddr string
	server    *server.MCPServer
}

// NewServer creates and initializes a new Server instance.
func NewServer(store *interceptor.TrafficStore, engine *rules.Engine, proxyAddr string) *Server {
	// Initialize the MCP server
	s := server.NewMCPServer("Agent Proxy", "1.0.0")

	ms := &Server{
		store:     store,
		engine:    engine,
		proxyAddr: proxyAddr,
		server:    s,
	}

	ms.registerTools()
	ms.registerResources()
	ms.registerPrompts()
	return ms
}

func (ms *Server) registerTools() {
	// Register list_traffic tool
	ms.server.AddTool(mcp.NewTool("list_traffic",
		mcp.WithDescription("List captured HTTP traffic summaries. Returns up to 20 recent entries."),
		mcp.WithString("filter", mcp.Description("Optional keyword to filter URL or Method")),
	), ms.listTrafficHandler)

	// Register get_traffic_details tool
	ms.server.AddTool(mcp.NewTool("get_traffic_details",
		mcp.WithDescription("Get full details of a specific traffic entry including headers and body."),
		mcp.WithString("id", mcp.Description("The ID of the traffic entry")),
	), ms.getTrafficDetailsHandler)

	// Register clear_traffic tool
	ms.server.AddTool(mcp.NewTool("clear_traffic",
		mcp.WithDescription("Clear all captured traffic logs."),
	), ms.clearTrafficHandler)

	// Register get_proxy_status tool
	ms.server.AddTool(mcp.NewTool("get_proxy_status",
		mcp.WithDescription("Get the current proxy address and status."),
	), ms.getProxyStatusHandler)

	// Register add_mock_rule tool
	ms.server.AddTool(mcp.NewTool("add_mock_rule",
		mcp.WithDescription("Add a mocking rule to intercept and return a static response for a specific URL."),
		mcp.WithString("url_pattern", mcp.Description("Keyword or pattern to match in URL")),
		mcp.WithString("method", mcp.Description("HTTP Method (e.g. GET, POST)")),
		mcp.WithNumber("status", mcp.Description("HTTP Status code to return (e.g. 200, 404)")),
		mcp.WithString("body", mcp.Description("Response body to return")),
	), ms.addMockRuleHandler)

	// --- New Tools ---

	// Rules Management
	ms.server.AddTool(mcp.NewTool("list_rules",
		mcp.WithDescription("List all active interception rules (mocks and breakpoints)."),
	), ms.listRulesHandler)

	ms.server.AddTool(mcp.NewTool("add_breakpoint_rule",
		mcp.WithDescription("Add a breakpoint rule to pause traffic for manual inspection."),
		mcp.WithString("url_pattern", mcp.Description("URL pattern to match")),
		mcp.WithString("method", mcp.Description("HTTP Method (optional)")),
		mcp.WithString("strategy", mcp.Description("Interception strategy: 'request', 'response', or 'both'")),
	), ms.addBreakpointRuleHandler)

	ms.server.AddTool(mcp.NewTool("delete_rule",
		mcp.WithDescription("Delete an interception rule by ID."),
		mcp.WithString("id", mcp.Description("The ID of the rule to delete")),
	), ms.deleteRuleHandler)

	// Traffic Execution
	ms.server.AddTool(mcp.NewTool("execute_request",
		mcp.WithDescription("Execute a custom HTTP request through the proxy. Can be used to replay an existing request by providing base_id."),
		mcp.WithString("method", mcp.Description("HTTP Method (e.g. GET, POST)")),
		mcp.WithString("url", mcp.Description("Target URL")),
		mcp.WithString("headers", mcp.Description("JSON string of headers (e.g. {\"Content-Type\": [\"application/json\"]})")),
		mcp.WithString("body", mcp.Description("Request body")),
		mcp.WithString("base_id", mcp.Description("Optional: The ID of an existing request to use as a template (replay)")),
	), ms.executeRequestHandler)
}

func (ms *Server) addMockRuleHandler(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing arguments")
	}

	urlPattern, _ := args["url_pattern"].(string)
	method, _ := args["method"].(string)
	status, _ := args["status"].(float64)
	body, _ := args["body"].(string)

	if urlPattern == "" {
		return nil, fmt.Errorf("url_pattern is required")
	}

	rule := &model.Rule{
		ID:         uuid.New().String(),
		Type:       model.RuleMock,
		URLPattern: urlPattern,
		Method:     method,
		Response: &model.MockResponse{
			Status: int(status),
			Body:   body,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"X-Mocked-By":  "Agent-Proxy",
			},
		},
	}

	ms.engine.AddRule(rule)
	return mcp.NewToolResultText(fmt.Sprintf("Mock rule added for %s %s (Returns %d)", method, urlPattern, int(status))), nil
}

func (ms *Server) listRulesHandler(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	rules := ms.engine.GetRules()
	var sb strings.Builder
	for _, r := range rules {
		sb.WriteString(fmt.Sprintf("ID: %s | Type: %s | Method: %s | Pattern: %s | Strategy: %s\n",
			r.ID, r.Type, r.Method, r.URLPattern, r.Strategy))
	}
	if sb.Len() == 0 {
		return mcp.NewToolResultText("No active rules."), nil
	}
	return mcp.NewToolResultText(sb.String()), nil
}

func (ms *Server) addBreakpointRuleHandler(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing arguments")
	}

	urlPattern, _ := args["url_pattern"].(string)
	method, _ := args["method"].(string)
	strategy, _ := args["strategy"].(string)

	if urlPattern == "" {
		return nil, fmt.Errorf("url_pattern is required")
	}

	rule := &model.Rule{
		ID:         uuid.New().String(),
		Type:       model.RuleBreakpoint,
		URLPattern: urlPattern,
		Method:     method,
		Strategy:   model.BreakpointStrategy(strategy),
	}

	ms.engine.AddRule(rule)
	return mcp.NewToolResultText(fmt.Sprintf("Breakpoint added for %s %s (Strategy: %s)", method, urlPattern, strategy)), nil
}

func (ms *Server) deleteRuleHandler(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing arguments")
	}
	id, _ := args["id"].(string)
	ms.engine.DeleteRule(id)
	return mcp.NewToolResultText(fmt.Sprintf("Rule %s deleted", id)), nil
}

func (ms *Server) executeRequestHandler(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing arguments")
	}

	method, _ := args["method"].(string)
	urlStr, _ := args["url"].(string)
	headersJSON, _ := args["headers"].(string)
	bodyStr, _ := args["body"].(string)
	baseID, _ := args["base_id"].(string)

	var finalMethod, finalURL, finalBody string
	finalHeaders := http.Header{}

	// 1. If base_id is provided, load the template
	if baseID != "" {
		entries, _ := ms.store.GetPage(0, 500)
		for _, e := range entries {
			if e.ID == baseID {
				finalMethod = e.Method
				finalURL = e.URL
				finalHeaders = e.RequestHeaders.Clone()
				finalBody = e.RequestBody
				break
			}
		}
	}

	// 2. Apply overrides from arguments
	if method != "" {
		finalMethod = method
	}
	if urlStr != "" {
		finalURL = urlStr
	}
	if bodyStr != "" {
		finalBody = bodyStr
	}
	if headersJSON != "" {
		var customHeaders map[string][]string
		if err := json.Unmarshal([]byte(headersJSON), &customHeaders); err == nil {
			for k, vs := range customHeaders {
				finalHeaders[k] = vs
			}
		}
	}

	if finalMethod == "" || finalURL == "" {
		return nil, fmt.Errorf("method and url are required (or valid base_id)")
	}

	req, err := http.NewRequest(finalMethod, finalURL, strings.NewReader(finalBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header = finalHeaders

	// Capture start
	entry, _ := interceptor.NewEntry(req)
	entry.ModifiedBy = "editor"

	client := &http.Client{Timeout: 30 * time.Second}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Capture response
	entry.Duration = time.Since(start)
	entry.Status = resp.StatusCode
	entry.ResponseHeaders = resp.Header.Clone()
	bodyBytes, _ := io.ReadAll(resp.Body)
	entry.ResponseBody = string(bodyBytes)

	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "image/") {
		encoded := base64.StdEncoding.EncodeToString(bodyBytes)
		entry.ResponseBody = fmt.Sprintf("data:%s;base64,%s", contentType, encoded)
	}

	ms.store.AddEntry(entry)

	return mcp.NewToolResultText(fmt.Sprintf("Request executed successfully.\nStatus: %d\nNew Entry ID: %s", resp.StatusCode, entry.ID)), nil
}

func (ms *Server) registerResources() {
	// Register proxy status resource
	ms.server.AddResource(mcp.NewResource("proxy://status",
		"Current Proxy Status",
		mcp.WithResourceDescription("Configuration and status of the agent proxy"),
		mcp.WithMIMEType("application/json"),
	), ms.proxyStatusResourceHandler)

	// Register latest traffic resource
	ms.server.AddResource(mcp.NewResource("traffic://latest",
		"Latest Traffic",
		mcp.WithResourceDescription("The most recent 10 HTTP requests captured"),
		mcp.WithMIMEType("application/json"),
	), ms.latestTrafficResourceHandler)
}

func (ms *Server) registerPrompts() {
	// Register analyze-traffic prompt
	ms.server.AddPrompt(mcp.NewPrompt("analyze-traffic",
		mcp.WithPromptDescription("Analyze recent traffic for errors or anomalies"),
	), ms.analyzeTrafficPromptHandler)

	// Register generate-api-docs prompt
	ms.server.AddPrompt(mcp.NewPrompt("generate-api-docs",
		mcp.WithPromptDescription("Generate API documentation from captured traffic"),
	), ms.generateAPIDocsPromptHandler)
}

func (ms *Server) analyzeTrafficPromptHandler(_ context.Context, _ mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	entries, _ := ms.store.GetPage(0, 5)

	var trafficData strings.Builder
	for _, e := range entries {
		trafficData.WriteString(fmt.Sprintf("[%s] %s (Status: %d)\n", e.Method, e.URL, e.Status))
	}

	return mcp.NewGetPromptResult("Analyze this traffic for errors:",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(trafficData.String())),
		},
	), nil
}

func (ms *Server) generateAPIDocsPromptHandler(_ context.Context, _ mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return mcp.NewGetPromptResult("Generate OpenAPI documentation based on the captured traffic logs. Focus on request/response structures and status codes.",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent("Please use the latest traffic logs to generate documentation.")),
		},
	), nil
}

func (ms *Server) listTrafficHandler(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		args = make(map[string]interface{})
	}
	filter, _ := args["filter"].(string)

	// Fetch a larger set to allow filtering, or we could add filtering to the repo
	entries, _ := ms.store.GetPage(0, 100)
	var results []string

	for _, e := range entries {
		if filter != "" {
			combined := strings.ToLower(e.Method + " " + e.URL)
			if !strings.Contains(combined, strings.ToLower(filter)) {
				continue
			}
		}

		line := fmt.Sprintf("[%s] %s (Status: %d, ID: %s)", e.Method, e.URL, e.Status, e.ID)
		results = append(results, line)

		if len(results) >= 20 {
			break
		}
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("No traffic found matching the criteria."), nil
	}

	return mcp.NewToolResultText(strings.Join(results, "\n")), nil
}

func (ms *Server) getTrafficDetailsHandler(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing arguments")
	}
	id, _ := args["id"].(string)
	if id == "" {
		return nil, fmt.Errorf("missing id")
	}

	// We search in the latest 100 entries
	entries, _ := ms.store.GetPage(0, 100)
	for _, e := range entries {
		if e.ID == id {
			details := fmt.Sprintf("ID: %s\nMethod: %s\nURL: %s\nStatus: %d\nDuration: %v\n\nRequest Headers:\n%v\n\nRequest Body:\n%s\n\nResponse Headers:\n%v\n\nResponse Body:\n%s",
				e.ID, e.Method, e.URL, e.Status, e.Duration, e.RequestHeaders, e.RequestBody, e.ResponseHeaders, e.ResponseBody)
			return mcp.NewToolResultText(details), nil
		}
	}

	return mcp.NewToolResultText("Traffic entry not found."), nil
}

func (ms *Server) clearTrafficHandler(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ms.store.ClearEntries()
	return mcp.NewToolResultText("Traffic logs cleared."), nil
}

func (ms *Server) getProxyStatusHandler(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	status := fmt.Sprintf("Proxy is running on: %s\nDashboard available on the API port", ms.proxyAddr)
	return mcp.NewToolResultText(status), nil
}

func (ms *Server) proxyStatusResourceHandler(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	rules := ms.engine.GetRules()
	status := fmt.Sprintf(`{"proxy_addr": "%s", "status": "running", "active_rules": %d}`, ms.proxyAddr, len(rules))
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "proxy://status",
			MIMEType: "application/json",
			Text:     status,
		},
	}, nil
}

func (ms *Server) latestTrafficResourceHandler(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	entries, _ := ms.store.GetPage(0, 10)

	var sb strings.Builder
	sb.WriteString("[\n")
	for i, e := range entries {
		sb.WriteString(fmt.Sprintf(`  {"id": "%s", "method": "%s", "url": "%s", "status": %d}`, e.ID, e.Method, e.URL, e.Status))
		if i < len(entries)-1 {
			sb.WriteString(",\n")
		}
	}
	sb.WriteString("\n]")

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "traffic://latest",
			MIMEType: "application/json",
			Text:     sb.String(),
		},
	}, nil
}

// StartSTDIO starts the MCP server using standard I/O.
func (ms *Server) StartSTDIO() error {
	return server.ServeStdio(ms.server)
}

// ServeSSE starts the MCP server using Server-Sent Events.
func (ms *Server) ServeSSE(addr string) error {
	// SSE server for MCP
	sse := server.NewSSEServer(ms.server)
	return sse.Start(addr)
}
