// Package db handles SQLite database initialization and schema management.
package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // SQLite driver
)

// DB is the global database connection.
var DB *sql.DB

// Init initializes the default database in the user's home directory.
func Init() {
	home, _ := os.UserHomeDir()
	dbPath := filepath.Join(home, ".glance.db")

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// High-performance SQLite settings for concurrent access
	DB.SetMaxOpenConns(1) // Force serialization to prevent "database is locked"

	if _, err := DB.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		log.Printf("Warning: Failed to enable WAL mode: %v", err)
	}
	if _, err := DB.Exec("PRAGMA synchronous=NORMAL;"); err != nil {
		log.Printf("Warning: Failed to set synchronous mode: %v", err)
	}

	createTables()
}

// InitCustom initializes the database at a specific path.
func InitCustom(path string) {
	var err error
	DB, err = sql.Open("sqlite", path)
	if err != nil {
		log.Fatalf("Failed to open database at %s: %v", path, err)
	}

	DB.SetMaxOpenConns(1)
	if _, err := DB.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		log.Printf("Warning: Failed to enable WAL mode: %v", err)
	}
	if _, err := DB.Exec("PRAGMA synchronous=NORMAL;"); err != nil {
		log.Printf("Warning: Failed to set synchronous mode: %v", err)
	}

	createTables()
}
func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS config (
			key TEXT PRIMARY KEY,
			value TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS traffic (
			id TEXT PRIMARY KEY,
			method TEXT,
			url TEXT,
			request_headers TEXT,
			request_body TEXT,
			response_headers TEXT,
			response_body TEXT,
			status INTEGER,
			start_time DATETIME,
			duration INTEGER,
			modified_by TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS rules (
			id TEXT PRIMARY KEY,
			type TEXT,
			url_pattern TEXT,
			method TEXT,
			strategy TEXT,
			response_json TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS scenarios (
			id TEXT PRIMARY KEY,
			name TEXT,
			description TEXT,
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS scenario_steps (
			id TEXT PRIMARY KEY,
			scenario_id TEXT,
			traffic_entry_id TEXT,
			step_order INTEGER,
			notes TEXT,
			FOREIGN KEY(scenario_id) REFERENCES scenarios(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS variable_mappings (
			id TEXT PRIMARY KEY,
			scenario_id TEXT,
			name TEXT,
			source_entry_id TEXT,
			source_path TEXT,
			target_json_path TEXT,
			FOREIGN KEY(scenario_id) REFERENCES scenarios(id) ON DELETE CASCADE
		)`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatalf("Failed to create table: %v\nQuery: %s", err, q)
		}
	}
}
