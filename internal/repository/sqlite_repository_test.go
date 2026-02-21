package repository

import (
	"database/sql"
	"glance/internal/model"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

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
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Fatalf("Failed to execute setup query: %v", err)
		}
	}
	return db
}

func TestSQLiteTrafficRepository(t *testing.T) {
	db := setupTestDB()
	repo := NewSQLiteTrafficRepository(db)

	entry := &model.TrafficEntry{
		ID:             "test-1",
		Method:         "GET",
		URL:            "http://test.local",
		RequestHeaders: http.Header{"X-Test": []string{"val"}},
		StartTime:      time.Now(),
	}

	if err := repo.Add(entry); err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	repo.Flush() // Ensure background write completes

	entries, total, err := repo.GetPage(0, 10)
	if err != nil || total != 1 {
		t.Errorf("GetPage failed: err=%v, total=%d", err, total)
	}

	if len(entries) != 1 || entries[0].ID != "test-1" {
		t.Errorf("Retrieved entry mismatch")
	}
}

func TestSQLiteConfigRepository(t *testing.T) {
	db := setupTestDB()
	repo := NewSQLiteConfigRepository(db)

	cfg := &model.Config{ProxyAddr: ":9999"}
	if err := repo.Save(cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	got, err := repo.Get()
	if err != nil || len(got.ProxyAddr) == 0 {
		t.Errorf("Get failed: err=%v, got=%+v", err, got)
	}
}

func TestSQLiteRuleRepository(t *testing.T) {
	db := setupTestDB()
	repo := NewSQLiteRuleRepository(db)

	rule := &model.Rule{
		ID:         "r1",
		Type:       model.RuleMock,
		URLPattern: "/api/test",
		Method:     "GET",
		Response: &model.MockResponse{
			Status: 200,
			Body:   "{\"ok\":true}",
		},
	}

	if err := repo.Add(rule); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	all, err := repo.GetAll()
	if err != nil || len(all) != 1 {
		t.Fatalf("GetAll failed: %v", err)
	}

	rule.URLPattern = "/api/updated"
	if err := repo.Update(rule); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	all, _ = repo.GetAll()
	if all[0].URLPattern != "/api/updated" {
		t.Errorf("Update not reflected")
	}

	if err := repo.Delete("r1"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	all, _ = repo.GetAll()
	if len(all) != 0 {
		t.Errorf("Delete failed")
	}
}

func TestSQLiteTrafficRepository_GetByIDs(t *testing.T) {
	db := setupTestDB()
	repo := NewSQLiteTrafficRepository(db)

	e1 := &model.TrafficEntry{ID: "t1", Method: "GET", URL: "u1", StartTime: time.Now()}
	e2 := &model.TrafficEntry{ID: "t2", Method: "POST", URL: "u2", StartTime: time.Now()}

	_ = repo.Add(e1)
	_ = repo.Add(e2)
	repo.Flush()

	got, err := repo.GetByIDs([]string{"t1", "t2"})
	if err != nil {
		t.Fatalf("GetByIDs failed: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(got))
	}
}

func TestSQLiteTrafficRepository_PruneAndClear(t *testing.T) {
	db := setupTestDB()
	repo := NewSQLiteTrafficRepository(db)

	for i := 0; i < 5; i++ {
		_ = repo.Add(&model.TrafficEntry{ID: string(rune(i)), StartTime: time.Now()})
	}
	repo.Flush()

	// Test Prune
	_ = repo.Prune(2)
	_, total, _ := repo.GetPage(0, 10)
	if total > 2 {
		t.Errorf("Expected max 2 entries, got %d", total)
	}

	// Test Clear
	_ = repo.Clear()
	_, total, _ = repo.GetPage(0, 10)
	if total != 0 {
		t.Errorf("Expected 0 entries, got %d", total)
	}
}

func TestSQLiteConfigRepository_Errors(t *testing.T) {
	db := setupTestDB()
	repo := NewSQLiteConfigRepository(db)

	// Test invalid JSON in DB
	_, _ = db.Exec("INSERT OR REPLACE INTO config (key, value) VALUES ('app_config', 'invalid json')")
	_, err := repo.Get()
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	// Test Get when no rows exist (should return default config)
	_, _ = db.Exec("DELETE FROM config")
	got, err := repo.Get()
	if err != nil || got == nil {
		t.Errorf("Expected default config when no rows exist, got err=%v", err)
	}
}

func TestSQLiteTrafficRepository_GetPage_Empty(t *testing.T) {
	db := setupTestDB()
	repo := NewSQLiteTrafficRepository(db)

	entries, total, err := repo.GetPage(0, 10)
	if err != nil || total != 0 || len(entries) != 0 {
		t.Errorf("Expected empty results, got err=%v, total=%d, len=%d", err, total, len(entries))
	}
}

func TestSQLiteTrafficRepository_GetByIDs_Empty(t *testing.T) {
	db := setupTestDB()
	repo := NewSQLiteTrafficRepository(db)

	got, err := repo.GetByIDs([]string{})
	if err != nil || len(got) != 0 {
		t.Errorf("Expected empty result for empty IDs, got err=%v, len=%d", err, len(got))
	}
}

func TestSQLiteScenarioRepository_Add_Error(t *testing.T) {
	db, _ := sql.Open("sqlite", ":memory:")
	repo := NewSQLiteScenarioRepository(db)
	_ = db.Close()
	err := repo.Add(&model.Scenario{ID: "s1"})
	if err == nil {
		t.Error("Expected error on closed DB in Add")
	}
}

func TestSQLiteScenarioRepository_Update_Error(t *testing.T) {
	db, _ := sql.Open("sqlite", ":memory:")
	repo := NewSQLiteScenarioRepository(db)
	_ = db.Close()
	err := repo.Update(&model.Scenario{ID: "s1"})
	if err == nil {
		t.Error("Expected error on closed DB in Update")
	}
}

func TestSQLiteTrafficRepository_GetPage_Error(t *testing.T) {
	db := setupTestDB()
	repo := NewSQLiteTrafficRepository(db)
	_ = db.Close()
	_, _, err := repo.GetPage(0, 10)
	if err == nil {
		t.Error("Expected error on closed DB in GetPage")
	}
}

func TestSQLiteRuleRepository_InvalidJSON(t *testing.T) {
	db := setupTestDB()
	repo := NewSQLiteRuleRepository(db)

	_, _ = db.Exec("INSERT INTO rules (id, type, url_pattern, method, strategy, response_json) VALUES ('r-bad', 'mock', '/bad', 'GET', 'request', 'invalid json')")
	all, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("Expected 1 rule, got %d", len(all))
	}
	if all[0].Response != nil {
		t.Errorf("Expected rule with nil response due to bad JSON, got %+v", all[0].Response)
	}
}

func TestSQLiteTrafficRepository_WriteError(t *testing.T) {
	db := setupTestDB()
	repo := NewSQLiteTrafficRepository(db)

	// Suppress expected error logs
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	entry := &model.TrafficEntry{ID: "dup", StartTime: time.Now()}
	_ = repo.Add(entry)
	repo.Flush()

	// Add same ID again - should trigger error in background worker
	_ = repo.Add(entry)
	repo.Flush()
}
