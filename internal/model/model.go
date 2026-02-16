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
