package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() {
	home, _ := os.UserHomeDir()
	dbPath := filepath.Join(home, ".agent-proxy.db")

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	createTables()
}

func InitCustom(path string) {
	var err error
	DB, err = sql.Open("sqlite", path)
	if err != nil {
		log.Fatalf("Failed to open database at %s: %v", path, err)
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
			duration INTEGER
		)`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatalf("Failed to create table: %v\nQuery: %s", err, q)
		}
	}
}
