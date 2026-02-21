package service

import (
	"testing"
)

func TestCAService(t *testing.T) {
	svc := NewCAService()
	cert := svc.GetCACert()
	if len(cert) == 0 {
		t.Error("Expected CA cert bytes, got empty")
	}
}
