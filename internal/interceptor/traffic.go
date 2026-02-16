// Package interceptor handles the capturing and temporary storage of HTTP traffic.
package interceptor

import (
	"agent-proxy/internal/config"
	"agent-proxy/internal/model"
	"agent-proxy/internal/repository"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TrafficStore provides an in-memory view and persistent storage for intercepted traffic.
type TrafficStore struct {
	repo repository.TrafficRepository
}

// NewTrafficStore creates a new TrafficStore with the provided repository.
func NewTrafficStore(repo repository.TrafficRepository) *TrafficStore {
	return &TrafficStore{repo: repo}
}

// AddEntry saves a new traffic entry to persistent storage.
func (s *TrafficStore) AddEntry(entry *model.TrafficEntry) {
	if s.repo == nil {
		return
	}

	cfg := config.Get()

	// 1. Enforce response size limit
	if cfg.MaxResponseSize > 0 && int64(len(entry.ResponseBody)) > cfg.MaxResponseSize {
		entry.ResponseBody = fmt.Sprintf("[Response body truncated. Size: %.2f MB exceeds limit of %.2f MB]",
			float64(len(entry.ResponseBody))/(1024*1024),
			float64(cfg.MaxResponseSize)/(1024*1024))
	}

	// 2. Save entry
	if err := s.repo.Add(entry); err != nil {
		log.Printf("Error saving traffic entry to repo: %v", err)
	}

	// 3. Auto-prune old history
	if cfg.HistoryLimit > 0 {
		if err := s.repo.Prune(cfg.HistoryLimit); err != nil {
			log.Printf("Error pruning history: %v", err)
		}
	}
}

// GetPage retrieves a paginated list of traffic entries.
func (s *TrafficStore) GetPage(offset, limit int) ([]*model.TrafficEntry, int) {
	if s.repo == nil {
		return nil, 0
	}
	entries, total, err := s.repo.GetPage(offset, limit)
	if err != nil {
		log.Printf("Error getting traffic page from repo: %v", err)
		return nil, 0
	}
	return entries, total
}

// ClearEntries removes all captured traffic from the repository.
func (s *TrafficStore) ClearEntries() {
	if s.repo == nil {
		return
	}
	if err := s.repo.Clear(); err != nil {
		log.Printf("Error clearing traffic in repo: %v", err)
	}
}

// ReadAndReplaceBody clones the request body without draining the original stream.
func ReadAndReplaceBody(r *http.Request) (string, error) {
	if r.Body == nil || r.Body == http.NoBody {
		return "", nil
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return string(body), nil
}

// ReadAndReplaceResponseBody clones the response body without draining the original stream.
func ReadAndReplaceResponseBody(res *http.Response) (string, error) {
	if res.Body == nil || res.Body == http.NoBody {
		return "", nil
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	res.Body = io.NopCloser(bytes.NewBuffer(body))

	contentType := res.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "image/") {
		encoded := base64.StdEncoding.EncodeToString(body)
		return fmt.Sprintf("data:%s;base64,%s", contentType, encoded), nil
	}

	return string(body), nil
}

// NewEntry creates a new TrafficEntry from an HTTP request.
func NewEntry(r *http.Request) (*model.TrafficEntry, error) {
	body, _ := ReadAndReplaceBody(r)
	return &model.TrafficEntry{
		ID:             uuid.New().String(),
		Method:         r.Method,
		URL:            r.URL.String(),
		RequestHeaders: r.Header.Clone(),
		RequestBody:    body,
		StartTime:      time.Now(),
	}, nil
}
