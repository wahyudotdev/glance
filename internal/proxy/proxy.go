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

	"agent-proxy/internal/ca"
	"agent-proxy/internal/interceptor"
	"agent-proxy/internal/model"
	"agent-proxy/internal/rules"
	"github.com/elazarl/goproxy"
)

type Breakpoint struct {
	ID      string
	Request *http.Request
	Entry   *model.TrafficEntry
	Resume  chan bool
	Abort   chan bool
}

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

func NewProxy(addr string) *Proxy {
	// Minimal fallback
	return NewProxyWithStore(addr, interceptor.NewTrafficStore(nil))
}

func NewProxyWithStore(addr string, store *interceptor.TrafficStore) *Proxy {
	ca.SetupCA()
	p := goproxy.NewProxyHttpServer()
	p.Verbose = false

	engine := rules.NewEngine()

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
	p.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			entry, err := interceptor.NewEntry(r)
			if err != nil {
				log.Printf("Error capturing request: %v", err)
			} else {
				ctx.UserData = entry
			}

			// Apply rules
			rule := engine.Match(r)
			if rule != nil {
				if rule.Type == rules.RuleMock && rule.Response != nil {
					resp := goproxy.NewResponse(r, goproxy.ContentTypeText, rule.Response.Status, rule.Response.Body)
					for k, v := range rule.Response.Headers {
						resp.Header.Set(k, v)
					}
					log.Printf("[MOCK] %s %s -> %d", r.Method, r.URL.String(), rule.Response.Status)
					return r, resp
				}

				if rule.Type == rules.RuleBreakpoint {
					log.Printf("[PAUSE] Intercepting %s %s", r.Method, r.URL.String())
					bp := &Breakpoint{
						ID:      entry.ID,
						Request: r,
						Entry:   entry,
						Resume:  make(chan bool),
						Abort:   make(chan bool),
					}

					proxy.bpMu.Lock()
					proxy.breakpoints[bp.ID] = bp
					proxy.bpMu.Unlock()

					if proxy.OnIntercept != nil {
						proxy.OnIntercept(bp)
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

					proxy.bpMu.Lock()
					delete(proxy.breakpoints, bp.ID)
					proxy.bpMu.Unlock()
				}
			}

			return r, nil
		},
	)

	// Capture Responses
	p.OnResponse().DoFunc(
		func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			if resp == nil {
				return resp
			}
			entry, ok := ctx.UserData.(*model.TrafficEntry)
			if ok && store != nil {
				body, _ := interceptor.ReadAndReplaceResponseBody(resp)
				entry.Status = resp.StatusCode
				entry.ResponseHeaders = resp.Header.Clone()
				entry.ResponseBody = body
				entry.Duration = time.Since(entry.StartTime)
				store.AddEntry(entry)

				if proxy.OnEntry != nil {
					proxy.OnEntry(entry)
				}

				log.Printf("[%d] %s %s (%v)", entry.Status, entry.Method, entry.URL, entry.Duration)
			}
			return resp
		},
	)

	return proxy
}

func (p *Proxy) Start() (string, error) {
	ln, err := net.Listen("tcp", p.addr)
	if err != nil {
		log.Printf("Port %s is in use, falling back to a random port...", p.addr)
		ln, err = net.Listen("tcp", ":0")
		if err != nil {
			return "", err
		}
	}

	actualAddr := ln.Addr().String()
	log.Printf("Proxy server starting on %s", actualAddr)

	// Use the listener with the server
	go http.Serve(ln, p.server)
	return actualAddr, nil
}

func (p *Proxy) GetBreakpoint(id string) *Breakpoint {
	p.bpMu.RLock()
	defer p.bpMu.RUnlock()
	return p.breakpoints[id]
}

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

func (p *Proxy) AbortRequest(id string) bool {
	bp := p.GetBreakpoint(id)
	if bp == nil {
		return false
	}
	bp.Abort <- true
	return true
}
