package mcp

import (
	"database/sql"
	glance_config "glance/internal/config"
	"glance/internal/interceptor"
	"glance/internal/model"
	"glance/internal/repository"
	"glance/internal/rules"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	_ "modernc.org/sqlite"
)

func setupTestServer() (*Server, *sql.DB, repository.TrafficRepository) {
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
	return ms, db, trafficRepo
}

func TestServer_StreamableHandler(t *testing.T) {
	ms, _, _ := setupTestServer()
	h := ms.GetStreamableHandler()
	if h == nil {
		t.Error("Expected non-nil handler")
	}
}

func TestToolHandlers(t *testing.T) {
	ms, _, repo := setupTestServer()

	t.Run("InspectNetworkTraffic", func(t *testing.T) {
		ms.store.AddEntry(&model.TrafficEntry{ID: "t1", Method: "GET", URL: "http://test.com"})
		ms.store.AddEntry(&model.TrafficEntry{ID: "t2", Method: "POST", URL: "http://api.com"})
		repo.Flush()

		// No filter
		res, _, _ := ms.handleInspectNetworkTraffic(listTrafficArgs{Limit: 10})
		if res == nil {
			t.Fatal("Expected result")
		}

		// With filter (matches)
		resF, _, _ := ms.handleInspectNetworkTraffic(listTrafficArgs{Filter: "test", Limit: 10})
		if resF == nil {
			t.Fatal("Expected filtered result")
		}

		// With filter (no match)
		resNM, _, _ := ms.handleInspectNetworkTraffic(listTrafficArgs{Filter: "nomatch", Limit: 10})
		if resNM == nil {
			t.Fatal("Expected no match result")
		}
	})

	t.Run("InspectRequestDetails", func(t *testing.T) {
		ms.store.AddEntry(&model.TrafficEntry{ID: "t2", Method: "POST", URL: "http://api.com"})
		res, _, err := ms.handleInspectRequestDetails(getTrafficDetailsArgs{ID: "t2"})
		if err != nil {
			t.Fatalf("Handle failed: %v", err)
		}
		if res == nil {
			t.Fatal("Expected result")
		}
	})

	t.Run("ClearTraffic", func(t *testing.T) {
		repo.Flush()
		_, _, _ = ms.handleClearTraffic()
		_, total := ms.store.GetPage(0, 10)
		if total != 0 {
			t.Errorf("Expected 0 entries, got %d", total)
		}
	})

	t.Run("ProxyStatus", func(t *testing.T) {
		res, _, _ := ms.handleGetProxyStatus()
		if res == nil {
			t.Error("Expected result")
		}
	})

	t.Run("MockRule", func(t *testing.T) {
		_, _, err := ms.handleAddMockRule(addMockRuleArgs{
			URLPattern: "test",
			Method:     "GET",
			Status:     200,
			Body:       "ok",
		})
		if err != nil {
			t.Fatalf("Handle failed: %v", err)
		}
		if len(ms.engine.GetRules()) != 1 {
			t.Error("Rule not added")
		}
	})

	t.Run("ListRules", func(t *testing.T) {
		res, _, _ := ms.handleListRules()
		if res == nil {
			t.Error("Expected result")
		}
	})

	t.Run("BreakpointRule", func(t *testing.T) {
		_, _, _ = ms.handleAddBreakpointRule(addBreakpointRuleArgs{
			URLPattern: "break",
			Strategy:   "both",
		})
		if len(ms.engine.GetRules()) != 2 {
			t.Error("Rule not added")
		}
	})

	t.Run("DeleteRule", func(t *testing.T) {
		rules := ms.engine.GetRules()
		_, _, _ = ms.handleDeleteRule(deleteRuleArgs{ID: rules[0].ID})
		if len(ms.engine.GetRules()) != 1 {
			t.Error("Rule not deleted")
		}
	})

	t.Run("ScenarioTools", func(t *testing.T) {
		// Add
		resA, _, errA := ms.handleAddScenario(addScenarioArgs{Name: "NewS", Description: "Desc"})
		if errA != nil {
			t.Fatalf("Add failed: %v", errA)
		}
		if resA == nil {
			t.Fatal("Expected result")
		}

		ms.store.AddEntry(&model.TrafficEntry{ID: "t_s1", Method: "GET", URL: "http://s1.com", ResponseBody: strings.Repeat("a", 2000)})
		repo.Flush()

		scenario := &model.Scenario{
			ID:   "s1",
			Name: "S1",
			Steps: []model.ScenarioStep{
				{TrafficEntryID: "t_s1", Order: 1, Notes: "Detailed Note"},
				{TrafficEntryID: "missing", Order: 2},
			},
			VariableMappings: []model.VariableMapping{
				{Name: "v1", SourceEntryID: "t_s1", SourcePath: "body.id", TargetJSONPath: "body.user_id"},
			},
		}
		_ = ms.scenarioRepo.Add(scenario)

		// List
		_, _, _ = ms.handleListScenarios()

		// Get
		res, _, err := ms.handleGetScenario(getScenarioArgs{ID: "s1"})
		if err != nil {
			t.Fatalf("Handle failed: %v", err)
		}
		if res == nil {
			t.Fatal("Expected result")
		}

		// Update with steps and mappings JSON
		_, _, _ = ms.handleUpdateScenario(updateScenarioArgs{
			ID:           "s1",
			Name:         "New Name",
			StepsJSON:    `[{"id":"step1", "traffic_entry_id":"t_s1", "order":1}]`,
			MappingsJSON: `[{"name":"v2", "source_entry_id":"t_s1", "source_path":"body.id", "target_json_path":"header.X"}]`,
		})

		// Delete
		_, _, _ = ms.handleDeleteScenario(deleteScenarioArgs{ID: "s1"})
	})

	t.Run("ExecuteRequest", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(200)
		}))
		defer ts.Close()

		// Basic execute
		_, _, err := ms.handleExecuteRequest(executeRequestArgs{
			Method: "GET",
			URL:    ts.URL,
		})
		if err != nil {
			t.Fatalf("Handle failed: %v", err)
		}

		// Execute with base_id
		ms.store.AddEntry(&model.TrafficEntry{ID: "base1", Method: "POST", URL: ts.URL})
		time.Sleep(50 * time.Millisecond)
		_, _, err = ms.handleExecuteRequest(executeRequestArgs{
			BaseID: "base1",
			Method: "GET", // override
		})
		if err != nil {
			t.Fatalf("Handle failed: %v", err)
		}

		// Missing required
		_, _, err = ms.handleExecuteRequest(executeRequestArgs{})
		if err == nil {
			t.Error("Expected error for missing method/url")
		}
	})

	t.Run("Resources", func(_ *testing.T) {
		_, _ = ms.handleReadProxyStatus(&mcp.ReadResourceRequest{})
		_, _ = ms.handleReadLatestTraffic(&mcp.ReadResourceRequest{})
	})

	t.Run("Prompts", func(_ *testing.T) {
		_, _ = ms.handlePromptAnalyzeTraffic(&mcp.GetPromptRequest{})
		_, _ = ms.handlePromptGenerateAPIDocs(&mcp.GetPromptRequest{})

		// Scenario test prompt
		s := &model.Scenario{ID: "s2", Name: "Test"}
		_ = ms.scenarioRepo.Add(s)
		_, _ = ms.handlePromptGenerateScenarioTest(&mcp.GetPromptRequest{
			Params: &mcp.GetPromptParams{
				Arguments: map[string]string{"id": "s2"},
			},
		})
	})
}

func TestRepositories(t *testing.T) {
	ms, _, _ := setupTestServer()

	t.Run("Traffic", func(t *testing.T) {
		entry := &model.TrafficEntry{ID: "t1", Method: "GET", URL: "http://test.com"}
		ms.store.AddEntry(entry)
		time.Sleep(100 * time.Millisecond)

		entries, total := ms.store.GetPage(0, 10)
		if total != 1 || entries[0].ID != "t1" {
			t.Errorf("Traffic not persisted correctly")
		}
	})

	t.Run("Rules", func(t *testing.T) {
		rule := &model.Rule{ID: "r1", Type: model.RuleMock, URLPattern: "/api"}
		ms.engine.AddRule(rule)
		rules := ms.engine.GetRules()
		if len(rules) != 1 || rules[0].ID != "r1" {
			t.Errorf("Rule not persisted")
		}

		ms.engine.DeleteRule("r1")
		if len(ms.engine.GetRules()) != 0 {
			t.Errorf("Rule not deleted")
		}
	})

	t.Run("Scenarios", func(t *testing.T) {
		scenario := &model.Scenario{ID: "s1", Name: "Flow 1"}
		if err := ms.scenarioRepo.Add(scenario); err != nil {
			t.Fatalf("Failed to add scenario: %v", err)
		}
		list, _ := ms.scenarioRepo.GetAll()
		if len(list) != 1 || list[0].Name != "Flow 1" {
			t.Errorf("Scenario not persisted")
		}

		if err := ms.scenarioRepo.Delete("s1"); err != nil {
			t.Errorf("Failed to delete scenario: %v", err)
		}
	})
}

func TestNewToolResultText(t *testing.T) {
	res := NewToolResultText("hello")
	if len(res.Content) != 1 {
		t.Fatal("Expected 1 content item")
	}
	// Content is an interface, we can use a type switch or cast if we know the implementation.
	// But the SDK makes it a bit hard to access private fields if any.
	// For now just checking it doesn't crash and has content.
}

func TestActiveSessions(t *testing.T) {
	ms, _, _ := setupTestServer()

	// Should be 0 initially
	if count := ms.ActiveSessions(); count != 0 {
		t.Errorf("Expected 0 sessions, got %d", count)
	}
}
