package mcp

import (
	"context"
	"database/sql"
	"encoding/json"
	glance_config "glance/internal/config"
	"glance/internal/interceptor"
	"glance/internal/model"
	"glance/internal/repository"
	"glance/internal/rules"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	_ "modernc.org/sqlite"
)

func setupTestServer() (*Server, *sql.DB) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(1)

	queries := []string{
		`CREATE TABLE config (key TEXT PRIMARY KEY, value TEXT)`,
		`CREATE TABLE traffic (
			id TEXT PRIMARY KEY, method TEXT, url TEXT,
			request_headers TEXT, request_body TEXT,
			response_headers TEXT, response_body TEXT,
			status INTEGER, start_time DATETIME, duration INTEGER, modified_by TEXT
		)`,
		`CREATE TABLE rules (
			id TEXT PRIMARY KEY, type TEXT, url_pattern TEXT,
			method TEXT, strategy TEXT, response_json TEXT
		)`,
		`CREATE TABLE scenarios (id TEXT PRIMARY KEY, name TEXT, description TEXT, created_at DATETIME)`,
		`CREATE TABLE scenario_steps (id TEXT PRIMARY KEY, scenario_id TEXT, traffic_entry_id TEXT, step_order INTEGER, notes TEXT)`,
		`CREATE TABLE variable_mappings (id TEXT PRIMARY KEY, scenario_id TEXT, name TEXT, source_entry_id TEXT, source_path TEXT, target_json_path TEXT)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Fatalf("Failed to execute setup query: %v", err)
		}
	}

	trafficRepo := repository.NewSQLiteTrafficRepository(db)
	ruleRepo := repository.NewSQLiteRuleRepository(db)
	scenarioRepo := repository.NewSQLiteScenarioRepository(db)
	configRepo := repository.NewSQLiteConfigRepository(db)

	glance_config.Init(configRepo)
	store := interceptor.NewTrafficStore(trafficRepo)
	engine := rules.NewEngine(ruleRepo)

	ms := NewServer(store, engine, ":8080", scenarioRepo)
	return ms, db
}

func TestTrafficHandlers(t *testing.T) {
	ms, _ := setupTestServer()
	ctx := context.Background()

	// 1. Add some traffic
	entry := &model.TrafficEntry{
		ID:     "t1",
		Method: "GET",
		URL:    "http://example.com/api",
		Status: 200,
	}
	ms.store.AddEntry(entry)
	time.Sleep(100 * time.Millisecond)
	ms.store.GetPage(0, 1) // Force repo to have data if async

	// 2. Test list_traffic
	t.Run("list_traffic", func(t *testing.T) {
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "list_traffic",
				Arguments: map[string]interface{}{
					"filter": "example",
				},
			},
		}
		res, err := ms.listTrafficHandler(ctx, req)
		if err != nil {
			t.Fatalf("list_traffic failed: %v", err)
		}
		if !strings.Contains(res.Content[0].(mcp.TextContent).Text, "http://example.com/api") {
			t.Errorf("Expected result to contain URL, got: %s", res.Content[0].(mcp.TextContent).Text)
		}
	})

	// 3. Test get_traffic_details
	t.Run("get_traffic_details", func(t *testing.T) {
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "get_traffic_details",
				Arguments: map[string]interface{}{
					"id": "t1",
				},
			},
		}
		res, err := ms.getTrafficDetailsHandler(ctx, req)
		if err != nil {
			t.Fatalf("get_traffic_details failed: %v", err)
		}
		if !strings.Contains(res.Content[0].(mcp.TextContent).Text, "ID: t1") {
			t.Errorf("Expected result to contain ID, got: %s", res.Content[0].(mcp.TextContent).Text)
		}
	})

	// 4. Test clear_traffic
	t.Run("clear_traffic", func(t *testing.T) {
		req := mcp.CallToolRequest{}
		_, err := ms.clearTrafficHandler(ctx, req)
		if err != nil {
			t.Fatalf("clear_traffic failed: %v", err)
		}
		// Verify empty
		entries, _ := ms.store.GetPage(0, 10)
		if len(entries) != 0 {
			t.Errorf("Expected 0 entries after clear, got %d", len(entries))
		}
	})
}

func TestRuleHandlers(t *testing.T) {
	ms, _ := setupTestServer()
	ctx := context.Background()

	t.Run("add_mock_rule", func(t *testing.T) {
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "add_mock_rule",
				Arguments: map[string]interface{}{
					"url_pattern": "/test",
					"method":      "GET",
					"status":      200.0,
					"body":        "hello",
				},
			},
		}
		_, err := ms.addMockRuleHandler(ctx, req)
		if err != nil {
			t.Fatalf("add_mock_rule failed: %v", err)
		}
		rules := ms.engine.GetRules()
		if len(rules) != 1 || rules[0].URLPattern != "/test" {
			t.Errorf("Rule not added correctly")
		}
	})

	t.Run("list_rules", func(t *testing.T) {
		req := mcp.CallToolRequest{}
		res, err := ms.listRulesHandler(ctx, req)
		if err != nil {
			t.Fatalf("list_rules failed: %v", err)
		}
		if !strings.Contains(res.Content[0].(mcp.TextContent).Text, "/test") {
			t.Errorf("Expected rule in list")
		}
	})

	t.Run("add_breakpoint_rule", func(t *testing.T) {
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "add_breakpoint_rule",
				Arguments: map[string]interface{}{
					"url_pattern": "/debug",
					"strategy":    "both",
				},
			},
		}
		_, err := ms.addBreakpointRuleHandler(ctx, req)
		if err != nil {
			t.Fatalf("add_breakpoint_rule failed: %v", err)
		}
		rules := ms.engine.GetRules()
		if len(rules) != 2 {
			t.Errorf("Expected 2 rules, got %d", len(rules))
		}
	})

	t.Run("delete_rule", func(t *testing.T) {
		rules := ms.engine.GetRules()
		id := rules[0].ID
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "delete_rule",
				Arguments: map[string]interface{}{
					"id": id,
				},
			},
		}
		_, err := ms.deleteRuleHandler(ctx, req)
		if err != nil {
			t.Fatalf("delete_rule failed: %v", err)
		}
		if len(ms.engine.GetRules()) != 1 {
			t.Errorf("Rule not deleted")
		}
	})
}

func TestScenarioHandlers(t *testing.T) {
	ms, _ := setupTestServer()
	ctx := context.Background()

	t.Run("add_scenario", func(t *testing.T) {
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "add_scenario",
				Arguments: map[string]interface{}{
					"name":        "Login Flow",
					"description": "Test login",
				},
			},
		}
		res, err := ms.addScenarioHandler(ctx, req)
		if err != nil {
			t.Fatalf("add_scenario failed: %v", err)
		}
		if !strings.Contains(res.Content[0].(mcp.TextContent).Text, "created with ID") {
			t.Errorf("Unexpected result: %s", res.Content[0].(mcp.TextContent).Text)
		}
	})

	t.Run("list_scenarios", func(t *testing.T) {
		res, err := ms.listScenariosHandler(ctx, mcp.CallToolRequest{})
		if err != nil {
			t.Fatalf("list_scenarios failed: %v", err)
		}
		if !strings.Contains(res.Content[0].(mcp.TextContent).Text, "Login Flow") {
			t.Errorf("Scenario not in list")
		}
	})

	t.Run("update_scenario", func(t *testing.T) {
		scenarios, _ := ms.scenarioRepo.GetAll()
		id := scenarios[0].ID

		steps := []model.ScenarioStep{{TrafficEntryID: "t1", Order: 1, Notes: "Step 1"}}
		stepsJSON, _ := json.Marshal(steps)

		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "update_scenario",
				Arguments: map[string]interface{}{
					"id":         id,
					"name":       "Updated Login",
					"steps_json": string(stepsJSON),
				},
			},
		}
		_, err := ms.updateScenarioHandler(ctx, req)
		if err != nil {
			t.Fatalf("update_scenario failed: %v", err)
		}

		updated, _ := ms.scenarioRepo.GetByID(id)
		if updated.Name != "Updated Login" || len(updated.Steps) != 1 {
			t.Errorf("Scenario not updated correctly: %+v", updated)
		}
	})

	t.Run("get_scenario", func(t *testing.T) {
		scenarios, _ := ms.scenarioRepo.GetAll()
		id := scenarios[0].ID
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "get_scenario",
				Arguments: map[string]interface{}{
					"id": id,
				},
			},
		}
		res, err := ms.getScenarioHandler(ctx, req)
		if err != nil {
			t.Fatalf("get_scenario failed: %v", err)
		}
		if !strings.Contains(res.Content[0].(mcp.TextContent).Text, "Updated Login") {
			t.Errorf("Scenario details missing")
		}
	})

	t.Run("delete_scenario", func(t *testing.T) {
		scenarios, _ := ms.scenarioRepo.GetAll()
		id := scenarios[0].ID
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "delete_scenario",
				Arguments: map[string]interface{}{
					"id": id,
				},
			},
		}
		_, err := ms.deleteScenarioHandler(ctx, req)
		if err != nil {
			t.Fatalf("delete_scenario failed: %v", err)
		}
		if list, _ := ms.scenarioRepo.GetAll(); len(list) != 0 {
			t.Errorf("Scenario not deleted")
		}
	})

	t.Run("execute_request", func(t *testing.T) {
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "execute_request",
				Arguments: map[string]interface{}{
					"method": "GET",
					"url":    "http://example.com/test-exec",
				},
			},
		}
		// Note: this will actually try to make an HTTP request if not mocked
		// For unit test, we might get an error but we check if it reaches the handler
		// But in setupTestServer, the handler is real.
		// Let's at least check it doesn't panic.
		_, _ = ms.executeRequestHandler(ctx, req)
	})

	t.Run("get_proxy_status", func(t *testing.T) {
		req := mcp.CallToolRequest{}
		res, err := ms.getProxyStatusHandler(ctx, req)
		if err != nil {
			t.Fatalf("get_proxy_status failed: %v", err)
		}
		if !strings.Contains(res.Content[0].(mcp.TextContent).Text, ":8080") {
			t.Errorf("Unexpected proxy status content")
		}
	})
}

func TestResourceHandlers(t *testing.T) {
	ms, _ := setupTestServer()

	t.Run("proxy_status_resource", func(t *testing.T) {
		res, err := ms.proxyStatusResourceHandler(context.Background(), mcp.ReadResourceRequest{})
		if err != nil {
			t.Fatalf("proxyStatusResourceHandler failed: %v", err)
		}
		if !strings.Contains(res[0].(mcp.TextResourceContents).Text, ":8080") {
			t.Errorf("Unexpected status: %s", res[0].(mcp.TextResourceContents).Text)
		}
	})

	t.Run("latest_traffic_resource", func(t *testing.T) {
		ms.store.AddEntry(&model.TrafficEntry{ID: "r1", Method: "POST", URL: "http://test.com"})
		time.Sleep(100 * time.Millisecond)
		res, err := ms.latestTrafficResourceHandler(context.Background(), mcp.ReadResourceRequest{})
		if err != nil {
			t.Fatalf("latestTrafficResourceHandler failed: %v", err)
		}
		if !strings.Contains(res[0].(mcp.TextResourceContents).Text, "r1") {
			t.Errorf("Unexpected traffic: %s", res[0].(mcp.TextResourceContents).Text)
		}
	})
}

func TestPromptHandlers(t *testing.T) {
	ms, _ := setupTestServer()

	t.Run("analyze_traffic_prompt", func(t *testing.T) {
		ms.store.AddEntry(&model.TrafficEntry{ID: "p1", Method: "GET", URL: "/error", Status: 500})
		time.Sleep(100 * time.Millisecond)
		res, err := ms.analyzeTrafficPromptHandler(context.Background(), mcp.GetPromptRequest{})
		if err != nil {
			t.Fatalf("analyzeTrafficPromptHandler failed: %v", err)
		}
		if !strings.Contains(res.Messages[0].Content.(mcp.TextContent).Text, "/error") {
			t.Errorf("Unexpected prompt content")
		}
	})

	t.Run("generate_scenario_test_prompt", func(t *testing.T) {
		scenario := &model.Scenario{ID: "s1", Name: "Test Flow"}
		ms.scenarioRepo.Add(scenario)

		req := mcp.GetPromptRequest{
			Params: mcp.GetPromptParams{
				Name: "generate-scenario-test",
				Arguments: map[string]string{
					"id": "s1",
				},
			},
		}
		res, err := ms.generateScenarioTestPromptHandler(context.Background(), req)
		if err != nil {
			t.Fatalf("generateScenarioTestPromptHandler failed: %v", err)
		}
		if !strings.Contains(res.Description, "Test Flow") {
			t.Errorf("Unexpected prompt description: %s", res.Description)
		}
	})
}
