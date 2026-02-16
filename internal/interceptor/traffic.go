package interceptor

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TrafficEntry struct {
	ID              string        `json:"id"`
	Method          string        `json:"method"`
	URL             string        `json:"url"`
	RequestHeaders  http.Header   `json:"request_headers"`
	RequestBody     string        `json:"request_body"`
	Status          int           `json:"status"`
	ResponseHeaders http.Header   `json:"response_headers"`
	ResponseBody    string        `json:"response_body"`
	StartTime       time.Time     `json:"start_time"`
	Duration        time.Duration `json:"duration"`
}

type TrafficStore struct {
	mu      sync.RWMutex
	entries []*TrafficEntry
}

func NewTrafficStore() *TrafficStore {
	return &TrafficStore{
		entries: make([]*TrafficEntry, 0),
	}
}

func (s *TrafficStore) AddEntry(entry *TrafficEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, entry)
}

func (s *TrafficStore) GetEntries() []*TrafficEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.entries
}

func (s *TrafficStore) ClearEntries() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = make([]*TrafficEntry, 0)
}

// Helper to clone request body without draining it
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

// Helper to clone response body without draining it
func ReadAndReplaceResponseBody(res *http.Response) (string, error) {
	if res.Body == nil || res.Body == http.NoBody {
		return "", nil
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	res.Body = io.NopCloser(bytes.NewBuffer(body))
	return string(body), nil
}

func NewEntry(r *http.Request) (*TrafficEntry, error) {
	body, _ := ReadAndReplaceBody(r)
	return &TrafficEntry{
		ID:             uuid.New().String(),
		Method:         r.Method,
		URL:            r.URL.String(),
		RequestHeaders: r.Header.Clone(),
		RequestBody:    body,
		StartTime:      time.Now(),
	}, nil
}
