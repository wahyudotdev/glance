package service

import (
	"glance/internal/interceptor"
	"glance/internal/model"
	"glance/internal/proxy"
	"glance/internal/rules"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/elazarl/goproxy"
)

func TestInterceptService(t *testing.T) {
	store := interceptor.NewTrafficStore(nil)
	engine := rules.NewEngine(&mockRuleRepo{rules: make(map[string]*model.Rule)})
	p := proxy.NewProxyWithRepositories(":0", store, engine)
	svc := NewInterceptService(p)

	t.Run("Abort - Missing ID", func(t *testing.T) {
		err := svc.Abort("non-existent")
		if err == nil {
			t.Error("Expected error")
		}
	})

	t.Run("Request Breakpoint and Abort", func(t *testing.T) {
		engine.AddRule(&model.Rule{ID: "b1", Enabled: true, Type: model.RuleBreakpoint, URLPattern: "abort", Strategy: "request"})
		req, _ := http.NewRequest("GET", "http://abort.me", nil)
		ctx := &goproxy.ProxyCtx{}

		done := make(chan bool)
		go func() {
			_, resp := p.HandleRequest(req, ctx)
			if resp != nil && resp.Body != nil {
				defer func() { _ = resp.Body.Close() }()
			}
			done <- true
		}()

		time.Sleep(50 * time.Millisecond)
		entry := ctx.UserData.(*model.TrafficEntry)

		if err := svc.Abort(entry.ID); err != nil {
			t.Errorf("Abort failed: %v", err)
		}
		<-done
	})

	t.Run("Response Breakpoint and Continue", func(t *testing.T) {
		engine.AddRule(&model.Rule{ID: "b2", Enabled: true, Type: model.RuleBreakpoint, URLPattern: "resume", Strategy: "response"})
		req, _ := http.NewRequest("GET", "http://resume.me", nil)
		res := &http.Response{
			StatusCode: 200,
			Request:    req,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
		}
		ctx := &goproxy.ProxyCtx{UserData: &model.TrafficEntry{ID: "t2"}}

		done := make(chan bool)
		go func() {
			resp := p.HandleResponse(res, ctx)
			if resp != nil && resp.Body != nil {
				defer func() { _ = resp.Body.Close() }()
			}
			done <- true
		}()

		time.Sleep(50 * time.Millisecond)
		if err := svc.ContinueResponse("t2", ContinueResponseParams{Status: 201}); err != nil {
			t.Errorf("ContinueResponse failed: %v", err)
		}
		<-done
	})

	t.Run("Request Breakpoint and Continue", func(t *testing.T) {
		engine.AddRule(&model.Rule{ID: "b3", Enabled: true, Type: model.RuleBreakpoint, URLPattern: "cont", Strategy: "request"})
		req, _ := http.NewRequest("GET", "http://cont.me", nil)
		ctx := &goproxy.ProxyCtx{}

		done := make(chan bool)
		go func() {
			_, resp := p.HandleRequest(req, ctx)
			if resp != nil && resp.Body != nil {
				defer func() { _ = resp.Body.Close() }()
			}
			done <- true
		}()

		time.Sleep(50 * time.Millisecond)
		entry := ctx.UserData.(*model.TrafficEntry)

		if err := svc.ContinueRequest(entry.ID, ContinueRequestParams{Method: "POST"}); err != nil {
			t.Errorf("ContinueRequest failed: %v", err)
		}
		<-done
	})

	t.Run("ContinueRequest - Missing ID", func(t *testing.T) {
		err := svc.ContinueRequest("non-existent", ContinueRequestParams{})
		if err == nil {
			t.Error("Expected error")
		}
	})

	t.Run("ContinueResponse - Missing ID", func(t *testing.T) {
		err := svc.ContinueResponse("non-existent", ContinueResponseParams{})
		if err == nil {
			t.Error("Expected error")
		}
	})
}
