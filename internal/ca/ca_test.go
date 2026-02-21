package ca

import (
	"os"
	"testing"
)

func TestSetupCA(t *testing.T) {
	// Call first time
	SetupCA()
	path1 := CAPath
	if path1 == "" {
		t.Error("Expected CAPath to be set")
	}
	if _, err := os.Stat(path1); os.IsNotExist(err) {
		t.Errorf("Expected CA file at %s to exist", path1)
	}

	// Call second time (idempotency/overwrite check)
	SetupCA()
	if CAPath != path1 {
		t.Errorf("CAPath changed on second call: %s -> %s", path1, CAPath)
	}
}

func TestSetupCA_WriteFailure(t *testing.T) {
	oldDir := getTmpDir
	defer func() { getTmpDir = oldDir }()

	// Mock non-existent/unwritable directory
	getTmpDir = func() string {
		return "/non-existent-dir-glance"
	}

	SetupCA()
	// Should not crash, just log warning
}
