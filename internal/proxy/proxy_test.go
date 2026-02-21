package proxy

import (
	"glance/internal/interceptor"
	"glance/internal/model"
	"glance/internal/rules"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/elazarl/goproxy"
)

type mockRuleRepo struct {
	rules []*model.Rule
}

func (m *mockRuleRepo) GetAll() ([]*model.Rule, error) { return m.rules, nil }
func (m *mockRuleRepo) Add(_ *model.Rule) error        { return nil }
func (m *mockRuleRepo) Update(_ *model.Rule) error     { return nil }
func (m *mockRuleRepo) Delete(_ string) error          { return nil }

func TestProxy_Start(t *testing.T) {
	p := NewProxy(":0")
	addr, err := p.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if addr == "" {
		t.Error("Expected non-empty address")
	}
}

func TestProxy_Constructors(t *testing.T) {
	p1 := NewProxy(":0")
	if p1 == nil || p1.addr != ":0" {
		t.Error("NewProxy failed")
	}

	p2 := NewProxyWithStore(":0", nil)
	if p2 == nil || p2.addr != ":0" {
		t.Error("NewProxyWithStore failed")
	}
}

func TestProxy_HandleRequest(t *testing.T) {
	repo := &mockRuleRepo{}
	engine := rules.NewEngine(repo)
	p := NewProxyWithRepositories(":0", nil, engine)

	t.Run("Normal Request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://test.com", nil)
		ctx := &goproxy.ProxyCtx{}
		r, resp := p.HandleRequest(req, ctx)
		if resp != nil {
			defer func() { _ = resp.Body.Close() }()
		}
		if r != req || resp != nil {
			t.Error("Expected original request and nil response")
		}
	})

	t.Run("CORS Preflight with Rule", func(t *testing.T) {
		repo.rules = []*model.Rule{{Method: "", URLPattern: "test.com"}}
		req, _ := http.NewRequest("OPTIONS", "http://test.com", nil)
		ctx := &goproxy.ProxyCtx{}
		_, resp := p.HandleRequest(req, ctx)
		if resp != nil {
			defer func() { _ = resp.Body.Close() }()
		}
		if resp == nil || resp.StatusCode != 204 {
			t.Errorf("Expected 204 response, got %v", resp)
		}
	})

	t.Run("Mock Response", func(t *testing.T) {
		repo.rules = []*model.Rule{{
			ID:         "m1",
			Type:       model.RuleMock,
			URLPattern: "mock.me",
			Response:   &model.MockResponse{Status: 201, Body: "mocked"},
		}}
		req, _ := http.NewRequest("GET", "http://mock.me", nil)
		ctx := &goproxy.ProxyCtx{}
		_, resp := p.HandleRequest(req, ctx)
		if resp != nil {
			defer func() { _ = resp.Body.Close() }()
		}
		if resp == nil || resp.StatusCode != 201 {
			t.Errorf("Expected 201 response, got %v", resp)
		}
	})

	t.Run("Breakpoint Request", func(t *testing.T) {
		repo.rules = []*model.Rule{{
			ID:         "b1",
			Type:       model.RuleBreakpoint,
			URLPattern: "pause.me",
			Strategy:   "request",
		}}
		req, _ := http.NewRequest("GET", "http://pause.me", nil)
		ctx := &goproxy.ProxyCtx{}

		// Run in goroutine because it blocks
		done := make(chan bool)
		go func() {
			_, resp := p.HandleRequest(req, ctx)
			if resp != nil {
				defer func() { _ = resp.Body.Close() }()
			}
			done <- true
		}()

		// Give it a moment to hit the breakpoint
		time.Sleep(50 * time.Millisecond)

		bp := p.GetBreakpoint(ctx.UserData.(*model.TrafficEntry).ID)
		if bp == nil {
			t.Fatal("Expected breakpoint to be registered")
		}

		p.ContinueRequest(bp.ID, "GET", "http://resumed.me", nil, "")
		<-done
	})

	t.Run("Abort Request", func(_ *testing.T) {
		repo.rules = []*model.Rule{{
			ID:         "b3",
			Type:       model.RuleBreakpoint,
			URLPattern: "abort.me",
			Strategy:   "request",
		}}
		req, _ := http.NewRequest("GET", "http://abort.me", nil)
		ctx := &goproxy.ProxyCtx{}

		done := make(chan bool)
		go func() {
			_, resp := p.HandleRequest(req, ctx)
			if resp != nil {
				defer func() { _ = resp.Body.Close() }()
				if resp.StatusCode == 502 {
					done <- true
				}
			}
		}()

		time.Sleep(50 * time.Millisecond)
		p.AbortRequest(ctx.UserData.(*model.TrafficEntry).ID)
		<-done
	})
}

func TestProxy_ContinueSuccess(t *testing.T) {
	p := NewProxyWithRepositories(":0", interceptor.NewTrafficStore(nil), rules.NewEngine(&mockRuleRepo{}))

	t.Run("ContinueRequest Success", func(t *testing.T) {
		entry := &model.TrafficEntry{ID: "c1"}
		req, _ := http.NewRequest("GET", "http://old.com", nil)
		p.bpMu.Lock()
		p.breakpoints["c1"] = &Breakpoint{ID: "c1", Request: req, Entry: entry, Resume: make(chan bool, 1)}
		p.bpMu.Unlock()

		success := p.ContinueRequest("c1", "POST", "http://new.com", nil, "body")
		if !success {
			t.Error("Expected success")
		}
		if req.Method != "POST" || req.URL.String() != "http://new.com" {
			t.Errorf("Request not updated: %s %s", req.Method, req.URL)
		}
	})

	t.Run("ContinueResponse Success", func(t *testing.T) {
		entry := &model.TrafficEntry{ID: "c2"}
		res := &http.Response{StatusCode: 200, Header: make(http.Header)}
		p.bpMu.Lock()
		p.breakpoints["c2"] = &Breakpoint{ID: "c2", Response: res, Entry: entry, Resume: make(chan bool, 1), Type: "response"}
		p.bpMu.Unlock()

		success := p.ContinueResponse("c2", 201, nil, "body")
		if !success {
			t.Error("Expected success")
		}
		if res.StatusCode != 201 {
			t.Errorf("Response not updated: %d", res.StatusCode)
		}
	})
}

func TestProxy_HandleResponse(t *testing.T) {
	store := interceptor.NewTrafficStore(nil)
	p := NewProxyWithRepositories(":0", store, rules.NewEngine(&mockRuleRepo{}))

	t.Run("Nil Response", func(t *testing.T) {
		res := p.HandleResponse(nil, nil)
		if res != nil {
			_ = res.Body.Close()
			t.Error("Expected nil")
		}
	})

	t.Run("Simple Response", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://test.com", nil)
		res := &http.Response{
			StatusCode: 200,
			Request:    req,
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     make(http.Header),
		}
		ctx := &goproxy.ProxyCtx{UserData: &model.TrafficEntry{}}
		got := p.HandleResponse(res, ctx)
		if got != nil && got.Body != nil {
			_ = got.Body.Close()
		}
		if got != res {
			t.Error("Expected original response")
		}
	})

	t.Run("Breakpoint Response", func(t *testing.T) {
		repo := &mockRuleRepo{rules: []*model.Rule{{
			ID:         "b2",
			Type:       model.RuleBreakpoint,
			URLPattern: "pause.res",
			Strategy:   "response",
		}}}
		store := interceptor.NewTrafficStore(nil)
		p := NewProxyWithRepositories(":0", store, rules.NewEngine(repo))

		req, _ := http.NewRequest("GET", "http://pause.res", nil)
		res := &http.Response{
			StatusCode: 200,
			Request:    req,
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     make(http.Header),
		}
		entry := &model.TrafficEntry{ID: "t2"}
		ctx := &goproxy.ProxyCtx{UserData: entry}

		done := make(chan bool)
		go func() {
			resp := p.HandleResponse(res, ctx)
			if resp != nil {
				defer func() { _ = resp.Body.Close() }()
			}
			done <- true
		}()

		time.Sleep(50 * time.Millisecond)
		bp := p.GetBreakpoint("t2")
		if bp == nil {
			t.Fatal("Expected breakpoint to be registered")
		}

		p.ContinueResponse("t2", 202, nil, "")
		<-done
	})
}

func TestProxy_AddBreakpointForTesting(t *testing.T) {
	p := NewProxy(":0")
	bp := &Breakpoint{ID: "test"}
	p.AddBreakpointForTesting(bp)
	if got := p.GetBreakpoint("test"); got != bp {
		t.Error("AddBreakpointForTesting failed")
	}
}

func TestProxy_BreakpointLogic(t *testing.T) {
	repo := &mockRuleRepo{}
	engine := rules.NewEngine(repo)
	store := interceptor.NewTrafficStore(nil)
	p := NewProxyWithRepositories(":0", store, engine)

	t.Run("GetBreakpoint - Empty", func(t *testing.T) {
		if bp := p.GetBreakpoint("123"); bp != nil {
			t.Error("Expected nil for non-existent breakpoint")
		}
	})

	t.Run("ContinueRequest - Missing ID", func(t *testing.T) {
		if success := p.ContinueRequest("123", "GET", "url", nil, ""); success {
			t.Error("Expected failure for non-existent breakpoint")
		}
	})

	t.Run("ContinueResponse - Missing ID", func(t *testing.T) {
		if success := p.ContinueResponse("123", 200, nil, ""); success {
			t.Error("Expected failure for non-existent breakpoint")
		}
	})

	t.Run("AbortRequest - Missing ID", func(t *testing.T) {
		if success := p.AbortRequest("123"); success {
			t.Error("Expected failure for non-existent breakpoint")
		}
	})
}
