package mcp

import (
	"database/sql"
	glance_config "glance/internal/config"
	"glance/internal/interceptor"
	"glance/internal/model"
	"glance/internal/repository"
	"glance/internal/rules"
	"log"
	"testing"
	"time"

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

// Since the new SDK uses a private 'handlers' map and generic AddTool,
// we can't easily call handlers directly by name without a session.
// However, we can test the repository-level logic which is what the handlers call.
// For a true integration test, we would use mcp.NewInMemoryTransports().

func TestRepositories(t *testing.T) {
	ms, _ := setupTestServer()

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
	})
}

func TestActiveSessions(t *testing.T) {
	ms, _ := setupTestServer()

	// Should be 0 initially
	if count := ms.ActiveSessions(); count != 0 {
		t.Errorf("Expected 0 sessions, got %d", count)
	}
}
