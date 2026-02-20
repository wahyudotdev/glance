// Package service implements the core business logic.
package service

import (
	"encoding/base64"
	"fmt"
	"glance/internal/interceptor"
	"glance/internal/model"
	"io"
	"net/http"
	"strings"
	"time"
)

// RequestService defines the interface for executing HTTP requests.
type RequestService interface {
	Execute(params ExecuteRequestParams) (*model.TrafficEntry, error)
}

// ExecuteRequestParams contains the parameters for executing a request.
type ExecuteRequestParams struct {
	Method  string
	URL     string
	Headers map[string][]string
	Body    string
}

type requestService struct {
	store *interceptor.TrafficStore
}

// NewRequestService creates a new RequestService.
func NewRequestService(store *interceptor.TrafficStore) RequestService {
	return &requestService{store: store}
}

func (s *requestService) Execute(params ExecuteRequestParams) (*model.TrafficEntry, error) {
	// Prepare the HTTP request
	req, err := http.NewRequest(params.Method, params.URL, strings.NewReader(params.Body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Copy headers
	for k, vs := range params.Headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	// Capture the request start
	entry, _ := interceptor.NewEntry(req)
	entry.ModifiedBy = "editor"

	// Execute
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	start := time.Now()
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	// Capture the response
	entry.Duration = time.Since(start)
	entry.Status = resp.StatusCode
	entry.ResponseHeaders = resp.Header.Clone()
	bodyBytes, _ := io.ReadAll(resp.Body)
	entry.ResponseBody = string(bodyBytes)

	// Detect and encode image if necessary (reuse logic from interceptor)
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "image/") {
		encoded := base64.StdEncoding.EncodeToString(bodyBytes)
		entry.ResponseBody = fmt.Sprintf("data:%s;base64,%s", contentType, encoded)
	}

	// Save to store
	s.store.AddEntry(entry)

	return entry, nil
}
