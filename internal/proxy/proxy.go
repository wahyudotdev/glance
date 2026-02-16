package proxy

import (
	"log"
	"net"
	"net/http"
	"time"

	"agent-proxy/internal/ca"
	"agent-proxy/internal/interceptor"
	"agent-proxy/internal/model"
	"agent-proxy/internal/rules"
	"github.com/elazarl/goproxy"
)

type Proxy struct {
	server  *goproxy.ProxyHttpServer
	addr    string
	Store   *interceptor.TrafficStore
	Engine  *rules.Engine
	OnEntry func(*model.TrafficEntry)
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
		server: p,
		addr:   addr,
		Store:  store,
		Engine: engine,
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
			if rule := engine.Match(r); rule != nil {
				if rule.Type == rules.RuleMock && rule.Response != nil {
					resp := goproxy.NewResponse(r, goproxy.ContentTypeText, rule.Response.Status, rule.Response.Body)
					for k, v := range rule.Response.Headers {
						resp.Header.Set(k, v)
					}
					log.Printf("[MOCK] %s %s -> %d", r.Method, r.URL.String(), rule.Response.Status)
					return r, resp
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
