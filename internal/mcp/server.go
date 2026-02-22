// Package mcp implements the Model Context Protocol server for AI agent integration.
package mcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"glance/internal/config"
	"glance/internal/interceptor"
	"glance/internal/model"
	"glance/internal/repository"
	"glance/internal/rules"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server manages the MCP connection and tool registrations.
type Server struct {
	store        *interceptor.TrafficStore
	engine       *rules.Engine
	scenarioRepo repository.ScenarioRepository
	proxyAddr    string
	server       *mcp.Server
}

type listTrafficArgs struct {
	Filter string  `json:"filter" jsonschema:"Optional keyword to filter URL or Method"`
	Limit  float64 `json:"limit" jsonschema:"Number of recent entries to return (default: 20)"`
}

type getTrafficDetailsArgs struct {
	ID string `json:"id" jsonschema:"The ID of the traffic entry"`
}

type addMockRuleArgs struct {
	URLPattern string  `json:"url_pattern" jsonschema:"Keyword or pattern to match in URL"`
	Method     string  `json:"method" jsonschema:"HTTP Method (e.g. GET, POST)"`
	Status     float64 `json:"status" jsonschema:"HTTP Status code to return (e.g. 200, 404)"`
	Body       string  `json:"body" jsonschema:"Response body to return"`
}

type addBreakpointRuleArgs struct {
	URLPattern string `json:"url_pattern" jsonschema:"URL pattern to match"`
	Method     string `json:"method" jsonschema:"HTTP Method (optional)"`
	Strategy   string `json:"strategy" jsonschema:"Interception strategy: 'request', 'response', or 'both'"`
}

type deleteRuleArgs struct {
	ID string `json:"id" jsonschema:"The ID of the rule to delete"`
}

type executeRequestArgs struct {
	Method  string `json:"method" jsonschema:"HTTP Method (e.g. GET, POST)"`
	URL     string `json:"url" jsonschema:"Target URL"`
	Headers string `json:"headers" jsonschema:"JSON string of headers (e.g. {\"Content-Type\": [\"application/json\"]})"`
	Body    string `json:"body" jsonschema:"Request body"`
	BaseID  string `json:"base_id" jsonschema:"Optional: The ID of an existing request to use as a template (replay)"`
}

type getScenarioArgs struct {
	ID string `json:"id" jsonschema:"The ID of the scenario"`
}

type addScenarioArgs struct {
	Name        string `json:"name" jsonschema:"The name of the scenario"`
	Description string `json:"description" jsonschema:"A brief description of the scenario"`
}

type updateScenarioArgs struct {
	ID           string `json:"id" jsonschema:"The ID of the scenario to update"`
	Name         string `json:"name" jsonschema:"New name (optional)"`
	Description  string `json:"description" jsonschema:"New description (optional)"`
	StepsJSON    string `json:"steps_json" jsonschema:"JSON array of steps"`
	MappingsJSON string `json:"mappings_json" jsonschema:"JSON array of mappings"`
}

type deleteScenarioArgs struct {
	ID string `json:"id" jsonschema:"The ID of the scenario to delete"`
}

// NewServer creates and initializes a new Server instance using the official SDK.
func NewServer(store *interceptor.TrafficStore, engine *rules.Engine, proxyAddr string, scenarioRepo repository.ScenarioRepository) *Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "Glance",
		Version: "1.0.0",
	}, nil)

	ms := &Server{
		store:        store,
		engine:       engine,
		scenarioRepo: scenarioRepo,
		proxyAddr:    proxyAddr,
		server:       s,
	}

	ms.registerTools()
	ms.registerResources()
	ms.registerPrompts()
	return ms
}

// ActiveSessions returns the number of currently active MCP sessions.
func (ms *Server) ActiveSessions() int {
	count := 0
	ms.server.Sessions()(func(_ *mcp.ServerSession) bool {
		count++
		return true
	})
	return count
}

// NewToolResultText is a helper to create a simple text-based tool result.
func NewToolResultText(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}
}

func (ms *Server) registerTools() {
	// 1. inspect_network_traffic
	mcp.AddTool(ms.server, &mcp.Tool{
		Name:        "inspect_network_traffic",
		Description: fmt.Sprintf("PRIMARY network debugging tool. MUST be called first to verify actual HTTP/HTTPS traffic when encountering errors, 4xx/5xx statuses, or unexpected API behavior. Returns a list of recent traffic summaries. Max limit follows system settings (currently %d).", config.Get().HistoryLimit),
	}, func(_ context.Context, _ *mcp.CallToolRequest, args listTrafficArgs) (*mcp.CallToolResult, any, error) {
		return ms.handleInspectNetworkTraffic(args)
	})

	// 2. inspect_request_details
	mcp.AddTool(ms.server, &mcp.Tool{
		Name:        "inspect_request_details",
		Description: "MANDATORY tool for deep inspection. Get full headers and body of a specific traffic entry to diagnose root causes of network failures.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args getTrafficDetailsArgs) (*mcp.CallToolResult, any, error) {
		return ms.handleInspectRequestDetails(args)
	})

	// 3. clear_traffic
	mcp.AddTool(ms.server, &mcp.Tool{
		Name:        "clear_traffic",
		Description: "Clear all captured traffic logs.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
		return ms.handleClearTraffic()
	})

	// 4. get_proxy_status
	mcp.AddTool(ms.server, &mcp.Tool{
		Name:        "get_proxy_status",
		Description: "Get the current proxy address and status.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
		return ms.handleGetProxyStatus()
	})

	// 5. add_mock_rule
	mcp.AddTool(ms.server, &mcp.Tool{
		Name:        "add_mock_rule",
		Description: "Add a mocking rule to intercept and return a static response for a specific URL.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args addMockRuleArgs) (*mcp.CallToolResult, any, error) {
		return ms.handleAddMockRule(args)
	})

	// 6. list_rules
	mcp.AddTool(ms.server, &mcp.Tool{
		Name:        "list_rules",
		Description: "List all active interception rules (mocks and breakpoints).",
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
		return ms.handleListRules()
	})

	// 7. add_breakpoint_rule
	mcp.AddTool(ms.server, &mcp.Tool{
		Name:        "add_breakpoint_rule",
		Description: "Add a breakpoint rule to pause traffic for manual inspection.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args addBreakpointRuleArgs) (*mcp.CallToolResult, any, error) {
		return ms.handleAddBreakpointRule(args)
	})

	// 8. delete_rule
	mcp.AddTool(ms.server, &mcp.Tool{
		Name:        "delete_rule",
		Description: "Delete an interception rule by ID.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args deleteRuleArgs) (*mcp.CallToolResult, any, error) {
		return ms.handleDeleteRule(args)
	})

	// 9. execute_request
	mcp.AddTool(ms.server, &mcp.Tool{
		Name:        "execute_request",
		Description: "Execute a custom HTTP request through the proxy. Can be used to replay an existing request by providing base_id.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args executeRequestArgs) (*mcp.CallToolResult, any, error) {
		return ms.handleExecuteRequest(args)
	})

	// 10. list_scenarios
	mcp.AddTool(ms.server, &mcp.Tool{
		Name:        "list_scenarios",
		Description: "List all recorded traffic scenarios (sequences of requests for test generation).",
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
		return ms.handleListScenarios()
	})

	// 11. get_scenario
	mcp.AddTool(ms.server, &mcp.Tool{
		Name:        "get_scenario",
		Description: "Get full details of a scenario, including the sequence of requests, responses, and variable mappings.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args getScenarioArgs) (*mcp.CallToolResult, any, error) {
		return ms.handleGetScenario(args)
	})

	// 12. add_scenario
	mcp.AddTool(ms.server, &mcp.Tool{
		Name:        "add_scenario",
		Description: "Create a new traffic scenario.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args addScenarioArgs) (*mcp.CallToolResult, any, error) {
		return ms.handleAddScenario(args)
	})

	// 13. update_scenario
	mcp.AddTool(ms.server, &mcp.Tool{
		Name:        "update_scenario",
		Description: "Update an existing scenario's metadata or steps.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args updateScenarioArgs) (*mcp.CallToolResult, any, error) {
		return ms.handleUpdateScenario(args)
	})

	// 14. delete_scenario
	mcp.AddTool(ms.server, &mcp.Tool{
		Name:        "delete_scenario",
		Description: "Delete a scenario by ID.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, args deleteScenarioArgs) (*mcp.CallToolResult, any, error) {
		return ms.handleDeleteScenario(args)
	})
}

func (ms *Server) handleInspectNetworkTraffic(args listTrafficArgs) (*mcp.CallToolResult, any, error) {
	limit := int(args.Limit)
	if limit <= 0 {
		limit = 20
	}

	cfg := config.Get()
	if limit > cfg.HistoryLimit && cfg.HistoryLimit > 0 {
		limit = cfg.HistoryLimit
	}

	entries, _ := ms.store.GetPage(0, limit)
	var results []string
	count := 0
	for _, e := range entries {
		if args.Filter != "" {
			combined := strings.ToLower(e.Method + " " + e.URL)
			if !strings.Contains(combined, strings.ToLower(args.Filter)) {
				continue
			}
		}
		line := fmt.Sprintf("[%s] %s (Status: %d, ID: %s)", e.Method, e.URL, e.Status, e.ID)
		results = append(results, line)
		count++
		if count >= limit {
			break
		}
	}
	if len(results) == 0 {
		return NewToolResultText("No traffic found matching the criteria."), nil, nil
	}
	return NewToolResultText(strings.Join(results, "\n")), nil, nil
}

func (ms *Server) handleInspectRequestDetails(args getTrafficDetailsArgs) (*mcp.CallToolResult, any, error) {
	cfg := config.Get()
	entries, _ := ms.store.GetPage(0, cfg.HistoryLimit)
	for _, e := range entries {
		if e.ID == args.ID {
			details := fmt.Sprintf("ID: %s\nMethod: %s\nURL: %s\nStatus: %d\nDuration: %v\n\nRequest Headers:\n%v\n\nRequest Body:\n%s\n\nResponse Headers:\n%v\n\nResponse Body:\n%s",
				e.ID, e.Method, e.URL, e.Status, e.Duration, e.RequestHeaders, e.RequestBody, e.ResponseHeaders, e.ResponseBody)
			return NewToolResultText(details), nil, nil
		}
	}
	return NewToolResultText("Traffic entry not found."), nil, nil
}

func (ms *Server) handleClearTraffic() (*mcp.CallToolResult, any, error) {
	ms.store.ClearEntries()
	return NewToolResultText("Traffic logs cleared."), nil, nil
}

func (ms *Server) handleGetProxyStatus() (*mcp.CallToolResult, any, error) {
	status := fmt.Sprintf("Proxy is running on: %s\nDashboard available on the API port", ms.proxyAddr)
	return NewToolResultText(status), nil, nil
}

func (ms *Server) handleAddMockRule(args addMockRuleArgs) (*mcp.CallToolResult, any, error) {
	rule := &model.Rule{
		ID:         uuid.New().String(),
		Enabled:    true,
		Type:       model.RuleMock,
		URLPattern: args.URLPattern,
		Method:     args.Method,
		Response: &model.MockResponse{
			Status: int(args.Status),
			Body:   args.Body,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"X-Mocked-By":  "Glance",
			},
		},
	}
	ms.engine.AddRule(rule)
	return NewToolResultText(fmt.Sprintf("Mock rule added for %s %s (Returns %d)", args.Method, args.URLPattern, int(args.Status))), nil, nil
}

func (ms *Server) handleListRules() (*mcp.CallToolResult, any, error) {
	rules := ms.engine.GetRules()
	var sb strings.Builder
	for _, r := range rules {
		status := "Enabled"
		if !r.Enabled {
			status = "Disabled"
		}
		fmt.Fprintf(&sb, "ID: %s | Status: %s | Type: %s | Method: %s | Pattern: %s | Strategy: %s\n",
			r.ID, status, r.Type, r.Method, r.URLPattern, r.Strategy)
	}
	if sb.Len() == 0 {
		return NewToolResultText("No active rules."), nil, nil
	}
	return NewToolResultText(sb.String()), nil, nil
}

func (ms *Server) handleAddBreakpointRule(args addBreakpointRuleArgs) (*mcp.CallToolResult, any, error) {
	rule := &model.Rule{
		ID:         uuid.New().String(),
		Enabled:    true,
		Type:       model.RuleBreakpoint,
		URLPattern: args.URLPattern,
		Method:     args.Method,
		Strategy:   model.BreakpointStrategy(args.Strategy),
	}
	ms.engine.AddRule(rule)
	return NewToolResultText(fmt.Sprintf("Breakpoint added for %s %s (Strategy: %s)", args.Method, args.URLPattern, args.Strategy)), nil, nil
}

func (ms *Server) handleDeleteRule(args deleteRuleArgs) (*mcp.CallToolResult, any, error) {
	ms.engine.DeleteRule(args.ID)
	return NewToolResultText(fmt.Sprintf("Rule %s deleted", args.ID)), nil, nil
}

func (ms *Server) handleExecuteRequest(args executeRequestArgs) (*mcp.CallToolResult, any, error) {
	var finalMethod, finalURL, finalBody string
	finalHeaders := http.Header{}

	if args.BaseID != "" {
		cfg := config.Get()
		entries, _ := ms.store.GetPage(0, cfg.HistoryLimit)
		for _, e := range entries {
			if e.ID == args.BaseID {
				finalMethod = e.Method
				finalURL = e.URL
				finalHeaders = e.RequestHeaders.Clone()
				finalBody = e.RequestBody
				break
			}
		}
	}

	if args.Method != "" {
		finalMethod = args.Method
	}
	if args.URL != "" {
		finalURL = args.URL
	}
	if args.Body != "" {
		finalBody = args.Body
	}
	if args.Headers != "" {
		var customHeaders map[string][]string
		if err := json.Unmarshal([]byte(args.Headers), &customHeaders); err == nil {
			for k, vs := range customHeaders {
				finalHeaders[k] = vs
			}
		}
	}

	if finalMethod == "" || finalURL == "" {
		return nil, nil, fmt.Errorf("method and url are required (or valid base_id)")
	}

	req, err := http.NewRequest(finalMethod, finalURL, strings.NewReader(finalBody))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header = finalHeaders

	entry, _ := interceptor.NewEntry(req)
	entry.ModifiedBy = "editor"

	client := &http.Client{Timeout: 30 * time.Second}
	start := time.Now()
	// #nosec G704 - This tool is intentionally designed to execute arbitrary requests as part of the MCP integration
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bodyBytes, _ := io.ReadAll(resp.Body)
	entry.Duration = time.Since(start)
	entry.Status = resp.StatusCode
	entry.ResponseHeaders = resp.Header.Clone()
	entry.ResponseBody = string(bodyBytes)

	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "image/") {
		encoded := base64.StdEncoding.EncodeToString(bodyBytes)
		entry.ResponseBody = fmt.Sprintf("data:%s;base64,%s", contentType, encoded)
	}

	ms.store.AddEntry(entry)
	return NewToolResultText(fmt.Sprintf("Request executed successfully.\nStatus: %d\nNew Entry ID: %s", resp.StatusCode, entry.ID)), nil, nil
}

func (ms *Server) handleListScenarios() (*mcp.CallToolResult, any, error) {
	scenarios, err := ms.scenarioRepo.GetAll()
	if err != nil {
		return nil, nil, err
	}
	var sb strings.Builder
	for _, s := range scenarios {
		fmt.Fprintf(&sb, "ID: %s | Name: %s | Description: %s\n", s.ID, s.Name, s.Description)
	}
	if sb.Len() == 0 {
		return NewToolResultText("No scenarios found."), nil, nil
	}
	return NewToolResultText(sb.String()), nil, nil
}

func (ms *Server) handleGetScenario(args getScenarioArgs) (*mcp.CallToolResult, any, error) {
	scenario, err := ms.scenarioRepo.GetByID(args.ID)
	if err != nil {
		return nil, nil, err
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Scenario: %s\nDescription: %s\n\n", scenario.Name, scenario.Description)
	sb.WriteString("Variable Mappings:\n")
	for _, m := range scenario.VariableMappings {
		fmt.Fprintf(&sb, "- Variable '%s' is extracted from %s (%s) and used in %s\n",
			m.Name, m.SourceEntryID, m.SourcePath, m.TargetJSONPath)
	}
	sb.WriteString("\nSteps (Sequence):\n")

	entries, _ := ms.store.GetPage(0, 1000)
	entryMap := make(map[string]*model.TrafficEntry)
	for _, e := range entries {
		entryMap[e.ID] = e
	}

	for _, step := range scenario.Steps {
		e, found := entryMap[step.TrafficEntryID]
		if !found {
			fmt.Fprintf(&sb, "%d. [MISSING ENTRY %s]\n", step.Order, step.TrafficEntryID)
			continue
		}
		fmt.Fprintf(&sb, "%d. [%s] %s (Status: %d)\n", step.Order, e.Method, e.URL, e.Status)
		if step.Notes != "" {
			fmt.Fprintf(&sb, "   Note: %s\n", step.Notes)
		}
		fmt.Fprintf(&sb, "   Request Headers: %v\n", e.RequestHeaders)
		if e.RequestBody != "" {
			fmt.Fprintf(&sb, "   Request Body: %s\n", e.RequestBody)
		}
		if e.ResponseBody != "" {
			body := e.ResponseBody
			if len(body) > 1000 {
				body = body[:1000] + "... [truncated]"
			}
			fmt.Fprintf(&sb, "   Response Body: %s\n", body)
		}
		sb.WriteString("\n")
	}
	return NewToolResultText(sb.String()), nil, nil
}

func (ms *Server) handleAddScenario(args addScenarioArgs) (*mcp.CallToolResult, any, error) {
	if args.Name == "" {
		return nil, nil, fmt.Errorf("name is required")
	}
	scenario := &model.Scenario{
		ID:          uuid.New().String(),
		Name:        args.Name,
		Description: args.Description,
		CreatedAt:   time.Now(),
	}
	if err := ms.scenarioRepo.Add(scenario); err != nil {
		return nil, nil, err
	}
	return NewToolResultText(fmt.Sprintf("Scenario '%s' created with ID: %s", args.Name, scenario.ID)), nil, nil
}

func (ms *Server) handleUpdateScenario(args updateScenarioArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return nil, nil, fmt.Errorf("id is required")
	}
	scenario, err := ms.scenarioRepo.GetByID(args.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("scenario not found: %v", err)
	}
	if args.Name != "" {
		scenario.Name = args.Name
	}
	if args.Description != "" {
		scenario.Description = args.Description
	}
	if args.StepsJSON != "" {
		var steps []model.ScenarioStep
		if err := json.Unmarshal([]byte(args.StepsJSON), &steps); err != nil {
			return nil, nil, fmt.Errorf("invalid steps_json: %v", err)
		}
		scenario.Steps = steps
	}
	if args.MappingsJSON != "" {
		var mappings []model.VariableMapping
		if err := json.Unmarshal([]byte(args.MappingsJSON), &mappings); err != nil {
			return nil, nil, fmt.Errorf("invalid mappings_json: %v", err)
		}
		scenario.VariableMappings = mappings
	}
	if err := ms.scenarioRepo.Update(scenario); err != nil {
		return nil, nil, err
	}
	return NewToolResultText(fmt.Sprintf("Scenario '%s' (%s) updated successfully.", scenario.Name, args.ID)), nil, nil
}

func (ms *Server) handleDeleteScenario(args deleteScenarioArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return nil, nil, fmt.Errorf("id is required")
	}
	if err := ms.scenarioRepo.Delete(args.ID); err != nil {
		return nil, nil, err
	}
	return NewToolResultText(fmt.Sprintf("Scenario %s deleted", args.ID)), nil, nil
}

func (ms *Server) registerResources() {
	ms.server.AddResource(&mcp.Resource{
		URI:      "proxy://status",
		Name:     "Current Proxy Status",
		MIMEType: "application/json",
	}, func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		return ms.handleReadProxyStatus(req)
	})

	ms.server.AddResource(&mcp.Resource{
		URI:      "traffic://latest",
		Name:     "Latest Traffic",
		MIMEType: "application/json",
	}, func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		return ms.handleReadLatestTraffic(req)
	})
}

func (ms *Server) handleReadProxyStatus(_ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	rules := ms.engine.GetRules()
	status := fmt.Sprintf(`{"proxy_addr": "%s", "status": "running", "active_rules": %d}`, ms.proxyAddr, len(rules))
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "proxy://status",
				MIMEType: "application/json",
				Text:     status,
			},
		},
	}, nil
}

func (ms *Server) handleReadLatestTraffic(_ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	entries, _ := ms.store.GetPage(0, 10)
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, e := range entries {
		fmt.Fprintf(&sb, `  {"id": "%s", "method": "%s", "url": "%s", "status": %d}`, e.ID, e.Method, e.URL, e.Status)
		if i < len(entries)-1 {
			sb.WriteString(",\n")
		}
	}
	sb.WriteString("\n]")
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "traffic://latest",
				MIMEType: "application/json",
				Text:     sb.String(),
			},
		},
	}, nil
}

func (ms *Server) registerPrompts() {
	ms.server.AddPrompt(&mcp.Prompt{
		Name:        "analyze-traffic",
		Description: "Analyze recent traffic for errors or anomalies",
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return ms.handlePromptAnalyzeTraffic(req)
	})

	ms.server.AddPrompt(&mcp.Prompt{
		Name:        "generate-api-docs",
		Description: "Generate API documentation from captured traffic",
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return ms.handlePromptGenerateAPIDocs(req)
	})

	ms.server.AddPrompt(&mcp.Prompt{
		Name:        "generate-scenario-test",
		Description: "Generate an automated test script for a specific scenario",
		Arguments: []*mcp.PromptArgument{
			{Name: "id", Description: "The ID of the scenario", Required: true},
			{Name: "framework", Description: "Test framework (playwright, cypress, go)"},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return ms.handlePromptGenerateScenarioTest(req)
	})
}

func (ms *Server) handlePromptAnalyzeTraffic(_ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	entries, _ := ms.store.GetPage(0, 5)
	var trafficData strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&trafficData, "[%s] %s (Status: %d)\n", e.Method, e.URL, e.Status)
	}
	return &mcp.GetPromptResult{
		Description: "Analyze this traffic for errors:",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: trafficData.String(),
				},
			},
		},
	}, nil
}

func (ms *Server) handlePromptGenerateAPIDocs(_ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Generate OpenAPI documentation based on the captured traffic logs.",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: "Please use the latest traffic logs to generate documentation.",
				},
			},
		},
	}, nil
}

func (ms *Server) handlePromptGenerateScenarioTest(req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	id := req.Params.Arguments["id"]
	framework := req.Params.Arguments["framework"]
	if framework == "" {
		framework = "playwright"
	}
	scenario, err := ms.scenarioRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	prompt := fmt.Sprintf("Please generate an automated test script using %s for the following scenario: '%s'.\n", framework, scenario.Name)
	prompt += "Use the variable mappings provided to handle dynamic values between requests."
	return &mcp.GetPromptResult{
		Description: prompt,
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: fmt.Sprintf("Scenario ID: %s", id),
				},
			},
		},
	}, nil
}

// StartSTDIO starts the MCP server using standard I/O.
func (ms *Server) StartSTDIO(ctx context.Context) error {
	return ms.server.Run(ctx, &mcp.StdioTransport{})
}

// ServeSSE starts the MCP server using Server-Sent Events.
func (ms *Server) ServeSSE(ctx context.Context, addr string) error {
	handler := ms.GetStreamableHandler()
	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		_ = server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}

// GetStreamableHandler returns an HTTP handler that supports Streamable HTTP (and SSE).
func (ms *Server) GetStreamableHandler() http.Handler {
	return mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		return ms.server
	}, nil)
}
