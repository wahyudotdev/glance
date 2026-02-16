package mcp

import (
	"context"
	"fmt"
	"strings"

	"agent-proxy/internal/interceptor"
	"agent-proxy/internal/rules"
	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	store     *interceptor.TrafficStore
	engine    *rules.Engine
	proxyAddr string
	server    *server.MCPServer
}

func NewMCPServer(store *interceptor.TrafficStore, engine *rules.Engine, proxyAddr string) *MCPServer {
	// Initialize the MCP server
	s := server.NewMCPServer("Agent Proxy", "1.0.0")

	ms := &MCPServer{
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

func (ms *MCPServer) registerTools() {
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
}

func (ms *MCPServer) addMockRuleHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	rule := &rules.Rule{
		ID:         uuid.New().String(),
		Type:       rules.RuleMock,
		URLPattern: urlPattern,
		Method:     method,
		Response: &rules.MockResponse{
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

func (ms *MCPServer) registerResources() {
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

func (ms *MCPServer) registerPrompts() {
	// Register analyze-traffic prompt
	ms.server.AddPrompt(mcp.NewPrompt("analyze-traffic",
		mcp.WithPromptDescription("Analyze recent traffic for errors or anomalies"),
	), ms.analyzeTrafficPromptHandler)

	// Register generate-api-docs prompt
	ms.server.AddPrompt(mcp.NewPrompt("generate-api-docs",
		mcp.WithPromptDescription("Generate API documentation from captured traffic"),
	), ms.generateAPIDocsPromptHandler)
}

func (ms *MCPServer) analyzeTrafficPromptHandler(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	entries := ms.store.GetEntries()
	count := 5
	if len(entries) < count {
		count = len(entries)
	}
	latest := entries[len(entries)-count:]

	var trafficData strings.Builder
	for _, e := range latest {
		trafficData.WriteString(fmt.Sprintf("[%s] %s (Status: %d)\n", e.Method, e.URL, e.Status))
	}

	return mcp.NewGetPromptResult("Analyze this traffic for errors:",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(trafficData.String())),
		},
	), nil
}

func (ms *MCPServer) generateAPIDocsPromptHandler(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return mcp.NewGetPromptResult("Generate OpenAPI documentation based on the captured traffic logs. Focus on request/response structures and status codes.",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent("Please use the latest traffic logs to generate documentation.")),
		},
	), nil
}

func (ms *MCPServer) listTrafficHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		args = make(map[string]interface{})
	}
	filter, _ := args["filter"].(string)

	entries := ms.store.GetEntries()
	var results []string

	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
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

func (ms *MCPServer) getTrafficDetailsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing arguments")
	}
	id, _ := args["id"].(string)
	if id == "" {
		return nil, fmt.Errorf("missing id")
	}

	entries := ms.store.GetEntries()
	for _, e := range entries {
		if e.ID == id {
			details := fmt.Sprintf("ID: %s\nMethod: %s\nURL: %s\nStatus: %d\nDuration: %v\n\nRequest Headers:\n%v\n\nRequest Body:\n%s\n\nResponse Headers:\n%v\n\nResponse Body:\n%s",
				e.ID, e.Method, e.URL, e.Status, e.Duration, e.RequestHeaders, e.RequestBody, e.ResponseHeaders, e.ResponseBody)
			return mcp.NewToolResultText(details), nil
		}
	}

	return mcp.NewToolResultText("Traffic entry not found."), nil
}

func (ms *MCPServer) clearTrafficHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ms.store.ClearEntries()
	return mcp.NewToolResultText("Traffic logs cleared."), nil
}

func (ms *MCPServer) getProxyStatusHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	status := fmt.Sprintf("Proxy is running on: %s\nDashboard available at: http://localhost:8081", ms.proxyAddr)
	return mcp.NewToolResultText(status), nil
}

func (ms *MCPServer) proxyStatusResourceHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	rules := ms.engine.GetRules()
	status := fmt.Sprintf(`{"proxy_addr": "%s", "dashboard_url": "http://localhost:8081", "status": "running", "active_rules": %d}`, ms.proxyAddr, len(rules))
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "proxy://status",
			MIMEType: "application/json",
			Text:     status,
		},
	}, nil
}

func (ms *MCPServer) latestTrafficResourceHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	entries := ms.store.GetEntries()
	count := 10
	if len(entries) < count {
		count = len(entries)
	}

	latest := entries[len(entries)-count:]
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, e := range latest {
		sb.WriteString(fmt.Sprintf(`  {"id": "%s", "method": "%s", "url": "%s", "status": %d}`, e.ID, e.Method, e.URL, e.Status))
		if i < len(latest)-1 {
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

func (ms *MCPServer) StartSTDIO() error {
	return server.ServeStdio(ms.server)
}

func (ms *MCPServer) ServeSSE(addr string) error {
	// SSE server for MCP
	sse := server.NewSSEServer(ms.server)
	return sse.Start(addr)
}
