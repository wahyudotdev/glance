package rules

import (
	"net/http"
	"strings"
	"sync"
)

type RuleType string

const (
	RuleMock       RuleType = "mock"
	RuleBreakpoint RuleType = "breakpoint"
)

type Rule struct {
	ID         string        `json:"id"`
	Type       RuleType      `json:"type"`
	URLPattern string        `json:"url_pattern"`
	Method     string        `json:"method"`
	Response   *MockResponse `json:"response,omitempty"`
}

type MockResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type Engine struct {
	mu    sync.RWMutex
	rules []*Rule
}

func NewEngine() *Engine {
	return &Engine{
		rules: make([]*Rule, 0),
	}
}

func (e *Engine) AddRule(rule *Rule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = append(e.rules, rule)
}

func (e *Engine) GetRules() []*Rule {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.rules
}

func (e *Engine) ClearRules() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = make([]*Rule, 0)
}

func (e *Engine) DeleteRule(id string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for i, r := range e.rules {
		if r.ID == id {
			e.rules = append(e.rules[:i], e.rules[i+1:]...)
			break
		}
	}
}

func (e *Engine) Match(r *http.Request) *Rule {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, rule := range e.rules {
		if rule.Method != "" && rule.Method != r.Method {
			continue
		}
		// Basic string contains for now, can be improved to regex
		if rule.URLPattern != "" && !strings.Contains(r.URL.String(), rule.URLPattern) {
			continue
		}
		return rule
	}
	return nil
}
