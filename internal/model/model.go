// Package model defines the data structures used throughout the Glance.
package model

import (
	"net/http"
	"time"
)

// TrafficEntry represents a single captured HTTP request/response pair.
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

// Config represents the application configuration.
type Config struct {
	ProxyAddr       string `json:"proxy_addr"`
	APIAddr         string `json:"api_addr"`
	MCPEnabled      bool   `json:"mcp_enabled"`
	HistoryLimit    int    `json:"history_limit"`
	MaxResponseSize int64  `json:"max_response_size"` // in bytes
	DefaultPageSize int    `json:"default_page_size"`
}

// RuleType defines the kind of interception rule.
type RuleType string

const (
	// RuleMock returns a static response.
	RuleMock RuleType = "mock"
	// RuleBreakpoint pauses the traffic.
	RuleBreakpoint RuleType = "breakpoint"
)

// BreakpointStrategy defines when to pause a request.
type BreakpointStrategy string

const (
	// StrategyRequest pauses before sending to target.
	StrategyRequest BreakpointStrategy = "request"
	// StrategyResponse pauses after receiving from target.
	StrategyResponse BreakpointStrategy = "response"
	// StrategyBoth pauses both before and after.
	StrategyBoth BreakpointStrategy = "both"
)

// Rule defines how to intercept specific traffic.
type Rule struct {
	ID         string             `json:"id"`
	Type       RuleType           `json:"type"`
	URLPattern string             `json:"url_pattern"`
	Method     string             `json:"method"`
	Strategy   BreakpointStrategy `json:"strategy,omitempty"` // For breakpoints
	Response   *MockResponse      `json:"response,omitempty"` // For mocks
}

// MockResponse defines the static response returned by a mock rule.
type MockResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// VariableMapping defines how a value from one response is used in a subsequent request.
type VariableMapping struct {
	Name           string `json:"name"`             // e.g., "sessionToken"
	SourceEntryID  string `json:"source_entry_id"`  // ID of the traffic entry providing the value
	SourcePath     string `json:"source_path"`      // e.g., "body.token" or "header.Set-Cookie"
	TargetJSONPath string `json:"target_json_path"` // e.g., "header.Authorization" or "body.user.id"
}

// ScenarioStep represents a single step in a recorded sequence.
type ScenarioStep struct {
	ID             string `json:"id"`
	TrafficEntryID string `json:"traffic_entry_id"`
	Order          int    `json:"order"`
	Notes          string `json:"notes,omitempty"`
}

// Scenario represents a sequence of related traffic entries for test generation.
type Scenario struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Steps            []ScenarioStep    `json:"steps"`
	VariableMappings []VariableMapping `json:"variable_mappings"`
	CreatedAt        time.Time         `json:"created_at"`
}
