package repository

import (
	"agent-proxy/internal/model"
	"database/sql"
	"log"
	"net/http"
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
		db.Exec(q)
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
	if err != nil || got.ProxyAddr != ":9999" {
		t.Errorf("Get failed: err=%v, got=%+v", err, got)
	}
}
