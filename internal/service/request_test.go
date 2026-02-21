package service

import (
	"glance/internal/interceptor"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestService_Execute(t *testing.T) {
	// Start a local server to test the execution
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer ts.Close()

	repo := &mockTrafficRepo{}
	store := interceptor.NewTrafficStore(repo)
	svc := NewRequestService(store)

	params := ExecuteRequestParams{
		Method: "GET",
		URL:    ts.URL,
		Headers: map[string][]string{
			"X-Custom": {"test"},
		},
		Body: "",
	}

	entry, err := svc.Execute(params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if entry.Status != 200 {
		t.Errorf("Expected status 200, got %d", entry.Status)
	}
	if entry.ResponseBody != `{"status":"ok"}` {
		t.Errorf("Expected body, got %s", entry.ResponseBody)
	}
	if len(repo.entries) != 1 {
		t.Errorf("Expected 1 entry in repo")
	}
}
