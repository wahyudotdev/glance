// Package proxy implements the core MITM proxy engine and interception logic.
package proxy

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"glance/internal/ca"
	"glance/internal/interceptor"
	"glance/internal/model"
	"glance/internal/rules"

	"github.com/elazarl/goproxy"
)

// Breakpoint represents a paused request or response waiting for user action.
type Breakpoint struct {
	ID       string
	Request  *http.Request
	Response *http.Response
	Entry    *model.TrafficEntry
	Resume   chan bool
	Abort    chan bool
	Type     string // "request" or "response"
}

// Proxy is the wrapper around the goproxy server that adds interception capabilities.
type Proxy struct {
	server      *goproxy.ProxyHttpServer
	addr        string
	Store       *interceptor.TrafficStore
	Engine      *rules.Engine
	OnEntry     func(*model.TrafficEntry)
	OnIntercept func(*Breakpoint) // Callback for UI notification

	breakpoints map[string]*Breakpoint
	bpMu        sync.RWMutex
}

// NewProxy creates a minimal Proxy instance.
func NewProxy(addr string) *Proxy {
	// Minimal fallback
	return NewProxyWithRepositories(addr, interceptor.NewTrafficStore(nil), rules.NewEngine(nil))
}

// NewProxyWithStore creates a Proxy with a custom traffic store.
func NewProxyWithStore(addr string, store *interceptor.TrafficStore) *Proxy {
	return NewProxyWithRepositories(addr, store, rules.NewEngine(nil))
}

// NewProxyWithRepositories creates a fully configured Proxy with store and rule engine.
func NewProxyWithRepositories(addr string, store *interceptor.TrafficStore, engine *rules.Engine) *Proxy {
	ca.SetupCA()
	p := goproxy.NewProxyHttpServer()
	p.Verbose = false

	proxy := &Proxy{
		server:      p,
		addr:        addr,
		Store:       store,
		Engine:      engine,
		breakpoints: make(map[string]*Breakpoint),
	}

	// Handle HTTPS CONNECT requests
	p.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	// Capture Requests and apply rules
	p.OnRequest().DoFunc(proxy.HandleRequest)

	// Capture Responses
	p.OnResponse().DoFunc(proxy.HandleResponse)

	return proxy
}

// HandleRequest processes an incoming HTTP request, applying rules and managing breakpoints.
func (p *Proxy) HandleRequest(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	entry, err := interceptor.NewEntry(r)
	if err != nil {
		log.Printf("Error capturing request: %v", err)
	} else {
		ctx.UserData = entry
	}

	// Apply rules
	rule := p.Engine.Match(r)

	// Handle CORS Preflight for any URL that has a rule
	if r.Method == "OPTIONS" && rule != nil {
		resp := goproxy.NewResponse(r, goproxy.ContentTypeText, 204, "")
		resp.Header.Set("Access-Control-Allow-Origin", "*")
		resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		resp.Header.Set("Access-Control-Allow-Headers", "*")
		resp.Header.Set("Access-Control-Max-Age", "86400")
		return r, resp
	}

	if rule != nil {
		if rule.Type == model.RuleMock && rule.Response != nil {
			entry.ModifiedBy = "mock"
			entry.Status = rule.Response.Status
			entry.ResponseHeaders = make(http.Header)
			for k, v := range rule.Response.Headers {
				entry.ResponseHeaders.Set(k, v)
			}
			entry.ResponseBody = rule.Response.Body
			entry.Duration = time.Since(entry.StartTime)

			// Save to store and broadcast
			if p.Store != nil {
				p.Store.AddEntry(entry)
			}
			if p.OnEntry != nil {
				p.OnEntry(entry)
			}

			resp := goproxy.NewResponse(r, goproxy.ContentTypeText, rule.Response.Status, rule.Response.Body)

			// Apply configured headers
			for k, v := range rule.Response.Headers {
				resp.Header.Set(k, v)
			}

			// Auto-inject CORS headers to prevent browser blocks
			resp.Header.Set("Access-Control-Allow-Origin", "*")
			resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			resp.Header.Set("Access-Control-Allow-Headers", "*")
			resp.Header.Set("Access-Control-Allow-Credentials", "true")

			log.Printf("[MOCK] %s %s -> %d", r.Method, r.URL.String(), rule.Response.Status)
			return r, resp
		}

		if rule.Type == model.RuleBreakpoint && (rule.Strategy == model.StrategyRequest || rule.Strategy == model.StrategyBoth || rule.Strategy == "") {
			entry.ModifiedBy = "breakpoint"
			log.Printf("[PAUSE REQ] Intercepting %s %s", r.Method, r.URL.String())
			bp := &Breakpoint{
				ID:      entry.ID,
				Request: r,
				Entry:   entry,
				Resume:  make(chan bool),
				Abort:   make(chan bool),
				Type:    "request",
			}
			p.bpMu.Lock()
			p.breakpoints[bp.ID] = bp
			p.bpMu.Unlock()

			if p.OnIntercept != nil {
				// Provide immediate feedback to UI and persist
				if p.Store != nil {
					p.Store.AddEntry(entry)
				}
				if p.OnEntry != nil {
					p.OnEntry(entry)
				}
				p.OnIntercept(bp)
			}

			// BLOCK here until resume or abort
			select {
			case <-bp.Resume:
				log.Printf("[RESUME] Resuming %s", bp.ID)
			case <-bp.Abort:
				log.Printf("[ABORT] Aborting %s", bp.ID)
				return r, goproxy.NewResponse(r, goproxy.ContentTypeText, 502, "Request aborted by user")
			case <-time.After(5 * time.Minute):
				log.Printf("[TIMEOUT] Auto-resuming %s after timeout", bp.ID)
			}

			p.bpMu.Lock()
			delete(p.breakpoints, bp.ID)
			p.bpMu.Unlock()
		}
	}

	return r, nil
}

// HandleResponse processes an outgoing HTTP response, capturing data and applying breakpoints.
func (p *Proxy) HandleResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if resp == nil {
		return resp
	}
	entry, ok := ctx.UserData.(*model.TrafficEntry)
	if ok && p.Store != nil {
		body, _ := interceptor.ReadAndReplaceResponseBody(resp)
		entry.Status = resp.StatusCode
		entry.ResponseHeaders = resp.Header.Clone()
		entry.ResponseBody = body
		entry.Duration = time.Since(entry.StartTime)

		// Check for Response Breakpoint
		rule := p.Engine.Match(resp.Request)
		if rule != nil && rule.Type == model.RuleBreakpoint && (rule.Strategy == model.StrategyResponse || rule.Strategy == model.StrategyBoth) {
			entry.ModifiedBy = "breakpoint"
			log.Printf("[PAUSE RES] Intercepting response for %s", resp.Request.URL.String())
			bp := &Breakpoint{
				ID:       entry.ID,
				Request:  resp.Request,
				Response: resp,
				Entry:    entry,
				Resume:   make(chan bool),
				Abort:    make(chan bool),
				Type:     "response",
			}

			p.bpMu.Lock()
			p.breakpoints[bp.ID] = bp
			p.bpMu.Unlock()

			if p.OnIntercept != nil {
				p.OnIntercept(bp)
			}

			// BLOCK here until resume or abort
			select {
			case <-bp.Resume:
				log.Printf("[RESUME RES] Resuming response for %s", bp.ID)
			case <-bp.Abort:
				log.Printf("[ABORT RES] Aborting response for %s", bp.ID)
				return goproxy.NewResponse(resp.Request, goproxy.ContentTypeText, 502, "Response aborted by user")
			case <-time.After(5 * time.Minute):
				log.Printf("[TIMEOUT RES] Auto-resuming response %s after timeout", bp.ID)
			}

			p.bpMu.Lock()
			delete(p.breakpoints, bp.ID)
			p.bpMu.Unlock()
		}

		p.Store.AddEntry(entry)

		if p.OnEntry != nil {
			p.OnEntry(entry)
		}

		log.Printf("[%d] %s %s (%v)", entry.Status, entry.Method, entry.URL, entry.Duration)
	}
	return resp
}

// Start begins the proxy server on the configured address.
func (p *Proxy) Start() (string, error) {
	ln, err := net.Listen("tcp", p.addr)
	if err != nil {
		log.Printf("Port %s is in use, falling back to a random port...", p.addr)
		ln, err = net.Listen("tcp", ":0") //nolint:gosec
		if err != nil {
			return "", err
		}
	}

	actualAddr := ln.Addr().String()

	// Use the listener with the server
	go http.Serve(ln, p.server) //nolint:errcheck,gosec
	return actualAddr, nil
}

// AddBreakpointForTesting is a helper for unit tests to register breakpoints manually.
func (p *Proxy) AddBreakpointForTesting(bp *Breakpoint) {
	p.bpMu.Lock()
	defer p.bpMu.Unlock()
	p.breakpoints[bp.ID] = bp
}

// GetBreakpoint retrieves an active breakpoint by its ID.
func (p *Proxy) GetBreakpoint(id string) *Breakpoint {
	p.bpMu.RLock()
	defer p.bpMu.RUnlock()
	return p.breakpoints[id]
}

// ContinueRequest resumes a paused request with potential modifications.
func (p *Proxy) ContinueRequest(id string, modifiedMethod, modifiedURL string, modifiedHeaders http.Header, modifiedBody string) bool {
	bp := p.GetBreakpoint(id)
	if bp == nil {
		return false
	}

	// Apply modifications to the original request
	if modifiedMethod != "" {
		bp.Request.Method = modifiedMethod
	}
	if modifiedURL != "" {
		newURL, err := url.Parse(modifiedURL)
		if err == nil {
			bp.Request.URL = newURL
		}
	}
	if modifiedHeaders != nil {
		bp.Request.Header = modifiedHeaders
	}
	if modifiedBody != "" {
		bp.Request.Body = io.NopCloser(strings.NewReader(modifiedBody))
		bp.Request.ContentLength = int64(len(modifiedBody))
	}

	// Update the entry for history consistency
	bp.Entry.Method = bp.Request.Method
	bp.Entry.URL = bp.Request.URL.String()
	bp.Entry.RequestHeaders = bp.Request.Header.Clone()
	bp.Entry.RequestBody = modifiedBody

	bp.Resume <- true
	return true
}

// AbortRequest terminates a paused request or response.
func (p *Proxy) AbortRequest(id string) bool {
	bp := p.GetBreakpoint(id)
	if bp == nil {
		return false
	}
	bp.Abort <- true
	return true
}

// ContinueResponse resumes a paused response with potential modifications.
func (p *Proxy) ContinueResponse(id string, modifiedStatus int, modifiedHeaders http.Header, modifiedBody string) bool {
	bp := p.GetBreakpoint(id)
	if bp == nil || bp.Type != "response" {
		return false
	}

	// Apply modifications to the original response
	if modifiedStatus > 0 {
		bp.Response.StatusCode = modifiedStatus
		bp.Response.Status = http.StatusText(modifiedStatus)
	}
	if modifiedHeaders != nil {
		bp.Response.Header = modifiedHeaders
	}
	if modifiedBody != "" {
		bp.Response.Body = io.NopCloser(strings.NewReader(modifiedBody))
		bp.Response.ContentLength = int64(len(modifiedBody))
	}

	// Update the entry for history consistency
	bp.Entry.Status = bp.Response.StatusCode
	bp.Entry.ResponseHeaders = bp.Response.Header.Clone()
	bp.Entry.ResponseBody = modifiedBody

	bp.Resume <- true
	return true
}
