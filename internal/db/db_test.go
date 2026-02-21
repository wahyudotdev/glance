package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCustom(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	InitCustom(dbPath)
	if DB == nil {
		t.Fatal("Expected DB to be initialized")
	}

	// Verify table existence
	_, err := DB.Exec("SELECT 1 FROM config LIMIT 1")
	if err != nil {
		t.Errorf("Config table not created: %v", err)
	}

	_ = DB.Close()
}

func TestInit_Default(t *testing.T) {
	// Mock HOME to a temp dir
	tmpHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer func() { _ = os.Setenv("HOME", oldHome) }()
	_ = os.Setenv("HOME", tmpHome)

	// Since Init might be called multiple times in different tests,
	// and it sets a global DB, we just ensure it doesn't panic.
	Init()
	if DB == nil {
		t.Error("Expected global DB to be set")
	}
}
