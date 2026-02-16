package proxy

import (
	"log"
	"net/http"
	"time"

	"agent-proxy/internal/ca"
	"agent-proxy/internal/interceptor"
	"github.com/elazarl/goproxy"
)

type Proxy struct {
	server *goproxy.ProxyHttpServer
	addr   string
	Store  *interceptor.TrafficStore
}

func NewProxy(addr string) *Proxy {
	ca.SetupCA()
	p := goproxy.NewProxyHttpServer()
	p.Verbose = false // Disable verbose to keep logs clean for now

	store := interceptor.NewTrafficStore()

	// Handle HTTPS CONNECT requests
	p.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	// Capture Requests
	p.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			entry, err := interceptor.NewEntry(r)
			if err != nil {
				log.Printf("Error capturing request: %v", err)
				return r, nil
			}
			ctx.UserData = entry
			return r, nil
		},
	)

	// Capture Responses
	p.OnResponse().DoFunc(
		func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			if resp == nil {
				return resp
			}
			entry, ok := ctx.UserData.(*interceptor.TrafficEntry)
			if ok {
				body, _ := interceptor.ReadAndReplaceResponseBody(resp)
				entry.Status = resp.StatusCode
				entry.ResponseHeaders = resp.Header.Clone()
				entry.ResponseBody = body
				entry.Duration = time.Since(entry.StartTime)
				store.AddEntry(entry)
				log.Printf("[%d] %s %s (%v)", entry.Status, entry.Method, entry.URL, entry.Duration)
			}
			return resp
		},
	)

	return &Proxy{
		server: p,
		addr:   addr,
		Store:  store,
	}
}

func (p *Proxy) Start() error {
	log.Printf("Proxy server starting on %s", p.addr)
	return http.ListenAndServe(p.addr, p.server)
}
