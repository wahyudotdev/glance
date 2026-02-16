package model

import (
	"net/http"
	"time"
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
	ModifiedBy      string        `json:"modified_by,omitempty"` // "mock" or "breakpoint"
}

type Config struct {
	ProxyAddr       string `json:"proxy_addr"`
	APIAddr         string `json:"api_addr"`
	MCPAddr         string `json:"mcp_addr"`
	MCPEnabled      bool   `json:"mcp_enabled"`
	HistoryLimit    int    `json:"history_limit"`
	MaxResponseSize int64  `json:"max_response_size"` // in bytes
	DefaultPageSize int    `json:"default_page_size"`
}

type RuleType string

const (
	RuleMock       RuleType = "mock"
	RuleBreakpoint RuleType = "breakpoint"
)

type BreakpointStrategy string

const (
	StrategyRequest  BreakpointStrategy = "request"
	StrategyResponse BreakpointStrategy = "response"
	StrategyBoth     BreakpointStrategy = "both"
)

type Rule struct {
	ID         string             `json:"id"`
	Type       RuleType           `json:"type"`
	URLPattern string             `json:"url_pattern"`
	Method     string             `json:"method"`
	Strategy   BreakpointStrategy `json:"strategy,omitempty"` // For breakpoints
	Response   *MockResponse      `json:"response,omitempty"` // For mocks
}

type MockResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}
