// Package service implements the core business logic.
package service

import (
	"fmt"
	"glance/internal/proxy"
	"net/http"
)

// InterceptService defines the interface for managing intercepted traffic.
type InterceptService interface {
	ContinueRequest(id string, params ContinueRequestParams) error
	ContinueResponse(id string, params ContinueResponseParams) error
	Abort(id string) error
}

// ContinueRequestParams contains parameters for resuming an intercepted request.
type ContinueRequestParams struct {
	Method  string
	URL     string
	Headers http.Header
	Body    string
}

// ContinueResponseParams contains parameters for resuming an intercepted response.
type ContinueResponseParams struct {
	Status  int
	Headers http.Header
	Body    string
}

type interceptService struct {
	proxy *proxy.Proxy
}

// NewInterceptService creates a new InterceptService.
func NewInterceptService(p *proxy.Proxy) InterceptService {
	return &interceptService{proxy: p}
}

func (s *interceptService) ContinueRequest(id string, params ContinueRequestParams) error {
	success := s.proxy.ContinueRequest(id, params.Method, params.URL, params.Headers, params.Body)
	if !success {
		return fmt.Errorf("intercepted request not found or already released")
	}
	return nil
}

func (s *interceptService) ContinueResponse(id string, params ContinueResponseParams) error {
	success := s.proxy.ContinueResponse(id, params.Status, params.Headers, params.Body)
	if !success {
		return fmt.Errorf("intercepted response not found or already released")
	}
	return nil
}

func (s *interceptService) Abort(id string) error {
	success := s.proxy.AbortRequest(id)
	if !success {
		return fmt.Errorf("intercepted request not found or already released")
	}
	return nil
}
