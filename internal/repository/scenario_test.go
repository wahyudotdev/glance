package repository

import (
	"database/sql"
	"glance/internal/model"
	"log"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func setupScenarioTestDB() *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	queries := []string{
		`CREATE TABLE scenarios (id TEXT PRIMARY KEY, name TEXT, description TEXT, created_at DATETIME)`,
		`CREATE TABLE traffic (id TEXT PRIMARY KEY, method TEXT, url TEXT, request_headers TEXT, request_body TEXT, response_headers TEXT, response_body TEXT, status INTEGER, start_time DATETIME, duration INTEGER, modified_by TEXT)`,
		`CREATE TABLE scenario_steps (id TEXT PRIMARY KEY, scenario_id TEXT, traffic_entry_id TEXT, step_order INTEGER, notes TEXT)`,
		`CREATE TABLE variable_mappings (id TEXT PRIMARY KEY, scenario_id TEXT, name TEXT, source_entry_id TEXT, source_path TEXT, target_json_path TEXT)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Fatalf("Failed to execute setup query: %v", err)
		}
	}
	return db
}

func TestSQLiteScenarioRepository(t *testing.T) {
	db := setupScenarioTestDB()
	repo := NewSQLiteScenarioRepository(db)

	scenario := &model.Scenario{
		ID:          "s1",
		Name:        "Test Scenario",
		Description: "A test sequence",
		CreatedAt:   time.Now(),
		Steps: []model.ScenarioStep{
			{ID: "step1", TrafficEntryID: "t1", Order: 1, Notes: "Note 1"},
			{ID: "step2", TrafficEntryID: "t2", Order: 2, Notes: "Note 2"},
		},
		VariableMappings: []model.VariableMapping{
			{Name: "token", SourceEntryID: "t1", SourcePath: "body.token", TargetJSONPath: "header.Auth"},
		},
	}

	if err := repo.Add(scenario); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Verify GetAll
	all, err := repo.GetAll()
	if err != nil || len(all) != 1 {
		t.Fatalf("GetAll failed: %v", err)
	}
	if all[0].Name != "Test Scenario" {
		t.Errorf("Name mismatch: %s", all[0].Name)
	}

	// Verify GetByID
	got, err := repo.GetByID("s1")
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if len(got.Steps) != 2 {
		t.Errorf("Step count mismatch: %d", len(got.Steps))
	}
	if len(got.VariableMappings) != 1 {
		t.Errorf("Mapping count mismatch: %d", len(got.VariableMappings))
	}

	// Test Update
	scenario.Name = "Updated Name"
	scenario.Steps = []model.ScenarioStep{
		{ID: "step3", TrafficEntryID: "t3", Order: 1},
	}
	if err := repo.Update(scenario); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	got, _ = repo.GetByID("s1")
	if got.Name != "Updated Name" || len(got.Steps) != 1 {
		t.Errorf("Update not applied correctly: %+v", got)
	}

	// Test Delete
	if err := repo.Delete("s1"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	all, _ = repo.GetAll()
	if len(all) != 0 {
		t.Errorf("Delete failed, still have %d scenarios", len(all))
	}
}
