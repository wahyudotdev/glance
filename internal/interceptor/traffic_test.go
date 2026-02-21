package interceptor

import (
	"bytes"
	"glance/internal/config"
	"glance/internal/model"
	"io"
	"net/http"
	"strings"
	"testing"
)

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

func TestReadAndReplaceResponseBody(t *testing.T) {
	content := "response content"
	res := &http.Response{
		Body:   io.NopCloser(bytes.NewBufferString(content)),
		Header: make(http.Header),
	}

	body, err := ReadAndReplaceResponseBody(res)
	if err != nil {
		t.Fatalf("ReadAndReplaceResponseBody failed: %v", err)
	}

	if body != content {
		t.Errorf("Expected body %q, got %q", content, body)
	}

	// Test with image
	imgContent := []byte{0x89, 0x50, 0x4E, 0x47} // PNG header
	res2 := &http.Response{
		Body:   io.NopCloser(bytes.NewBuffer(imgContent)),
		Header: http.Header{"Content-Type": []string{"image/png"}},
	}
	body2, _ := ReadAndReplaceResponseBody(res2)
	if body2 == string(imgContent) {
		t.Errorf("Expected base64 encoded body for image")
	}
}

type mockRepo struct {
	entries []*model.TrafficEntry
}

func (m *mockRepo) Add(e *model.TrafficEntry) error {
	m.entries = append(m.entries, e)
	return nil
}
func (m *mockRepo) GetPage(_, _ int) ([]*model.TrafficEntry, int, error) {
	return m.entries, len(m.entries), nil
}
func (m *mockRepo) GetByIDs(_ []string) ([]*model.TrafficEntry, error) { return nil, nil }
func (m *mockRepo) Clear() error                                       { return nil }
func (m *mockRepo) Prune(_ int) error                                  { return nil }
func (m *mockRepo) Flush()                                             {}

type mockConfigRepo struct{}

func (m *mockConfigRepo) Get() (*model.Config, error) { return &model.Config{HistoryLimit: 100}, nil }
func (m *mockConfigRepo) Save(_ *model.Config) error  { return nil }

func TestTrafficStore(t *testing.T) {
	config.Init(&mockConfigRepo{})
	repo := &mockRepo{}
	store := NewTrafficStore(repo)

	entry := &model.TrafficEntry{ID: "test-1"}
	store.AddEntry(entry)

	if len(repo.entries) != 1 {
		t.Errorf("Expected 1 entry in repo")
	}

	entries, total := store.GetPage(0, 10)
	if total != 1 || entries[0].ID != "test-1" {
		t.Errorf("GetPage failed")
	}

	store.ClearEntries()
}

func TestTrafficStore_AddEntry_Truncation(t *testing.T) {
	// Set a very small limit
	cfg := &model.Config{MaxResponseSize: 10, HistoryLimit: 100}
	config.Init(&mockConfigRepoForTrunc{cfg: cfg})

	repo := &mockRepo{}
	store := NewTrafficStore(repo)

	entry := &model.TrafficEntry{
		ID:           "test-2",
		ResponseBody: "this is a very long response body that should be truncated",
	}
	store.AddEntry(entry)

	if !strings.Contains(entry.ResponseBody, "truncated") {
		t.Errorf("Expected truncated message, got %s", entry.ResponseBody)
	}
}

type mockConfigRepoForTrunc struct {
	cfg *model.Config
}

func (m *mockConfigRepoForTrunc) Get() (*model.Config, error) { return m.cfg, nil }
func (m *mockConfigRepoForTrunc) Save(_ *model.Config) error  { return nil }

func TestTrafficStore_NoRepo(_ *testing.T) {
	store := NewTrafficStore(nil)
	store.AddEntry(&model.TrafficEntry{ID: "1"})
	store.GetPage(0, 10)
	store.ClearEntries()
}
