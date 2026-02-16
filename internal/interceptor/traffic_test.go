package interceptor

import (
	"agent-proxy/internal/model"
	"bytes"
	"io"
	"net/http"
	"testing"
)

type mockTrafficRepo struct {
	entries []*model.TrafficEntry
}

func (m *mockTrafficRepo) Add(e *model.TrafficEntry) error { m.entries = append(m.entries, e); return nil }
func (m *mockTrafficRepo) GetPage(o, l int) ([]*model.TrafficEntry, int, error) {
	return m.entries, len(m.entries), nil
}
func (m *mockTrafficRepo) Clear() error        { m.entries = nil; return nil }
func (m *mockTrafficRepo) Prune(l int) error   { return nil }

func TestReadAndReplaceBody(t *testing.T) {
	content := "test body content"
	req, _ := http.NewRequest("POST", "http://test.com", bytes.NewBufferString(content))

	body, err := ReadAndReplaceBody(req)
	if err != nil {
		t.Fatalf("ReadAndReplaceBody failed: %v", err)
	}

	if body != content {
		t.Errorf("Expected body %q, got %q", content, body)
	}

	// Verify body is still readable (was replaced)
	newBody, _ := io.ReadAll(req.Body)
	if string(newBody) != content {
		t.Errorf("Request body was drained or corrupted")
	}
}

func TestNewEntry(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/path?query=1", nil)
	req.Header.Set("X-Test", "Value")

	entry, err := NewEntry(req)
	if err != nil {
		t.Fatalf("NewEntry failed: %v", err)
	}

	if entry.Method != "GET" {
		t.Errorf("Expected GET, got %s", entry.Method)
	}
	if entry.URL != "http://example.com/path?query=1" {
		t.Errorf("URL mismatch: %s", entry.URL)
	}
	if entry.RequestHeaders.Get("X-Test") != "Value" {
		t.Errorf("Header not captured")
	}
}
