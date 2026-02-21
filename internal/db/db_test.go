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

	// Support both HOME (Unix) and USERPROFILE (Windows) for local runs
	oldHome := os.Getenv("HOME")
	oldUP := os.Getenv("USERPROFILE")
	defer func() {
		_ = os.Setenv("HOME", oldHome)
		_ = os.Setenv("USERPROFILE", oldUP)
	}()
	_ = os.Setenv("HOME", tmpHome)
	_ = os.Setenv("USERPROFILE", tmpHome)

	// Since Init might be called multiple times in different tests,
	// and it sets a global DB, we just ensure it doesn't panic.
	Init()
	if DB == nil {
		t.Error("Expected global DB to be set")
	}

	// Ensure tables exist after default Init
	_, err := DB.Exec("SELECT 1 FROM scenarios LIMIT 1")
	if err != nil {
		t.Errorf("Tables not created in default Init: %v", err)
	}
}

func TestInitCustom_Twice(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test2.db")

	// First init
	InitCustom(dbPath)
	db1 := DB

	// Second init (re-open)
	InitCustom(dbPath)
	if DB == nil {
		t.Error("Expected DB to be initialized")
	}
	if DB == db1 {
		// This might happen if sql.Open returns the same instance, but usually it's different.
		// We just want to ensure it's not nil and maybe log it.
		t.Log("DB instance is the same after re-init")
	}
}

func TestInit_Failures(t *testing.T) {
	oldFatalf := fatalf
	defer func() { fatalf = oldFatalf }()

	var fatalCalled bool
	fatalf = func(_ string, _ ...any) {
		fatalCalled = true
	}

	// Test InitCustom with invalid path
	InitCustom("/non-existent/dir/db.db")
	if !fatalCalled {
		t.Error("Expected fatalf to be called for invalid path")
	}

	// Reset
	fatalCalled = false

	// Test createTables with closed DB
	if DB != nil {
		_ = DB.Close()
		createTables()
		if !fatalCalled {
			t.Error("Expected fatalf to be called for closed DB in createTables")
		}
	}
}
