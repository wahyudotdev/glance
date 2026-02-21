package ca

import (
	"os"
	"testing"
)

func TestSetupCA(t *testing.T) {
	SetupCA()
	if CAPath == "" {
		t.Error("Expected CAPath to be set")
	}
	if _, err := os.Stat(CAPath); os.IsNotExist(err) {
		t.Errorf("Expected CA file at %s to exist", CAPath)
	}
}
