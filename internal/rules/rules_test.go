package rules

import (
	"errors"
	"glance/internal/model"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

type mockRuleRepo struct {
	rules []*model.Rule
	err   error
}

func (m *mockRuleRepo) GetAll() ([]*model.Rule, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rules, nil
}
func (m *mockRuleRepo) Add(r *model.Rule) error {
	if m.err != nil {
		return m.err
	}
	m.rules = append(m.rules, r)
	return nil
}
func (m *mockRuleRepo) Update(r *model.Rule) error {
	if m.err != nil {
		return m.err
	}
	for i, existing := range m.rules {
		if existing.ID == r.ID {
			m.rules[i] = r
			break
		}
	}
	return nil
}
func (m *mockRuleRepo) Delete(id string) error {
	if m.err != nil {
		return m.err
	}
	for i, r := range m.rules {
		if r.ID == id {
			m.rules = append(m.rules[:i], m.rules[i+1:]...)
			break
		}
	}
	return nil
}

func TestEngine_Match(t *testing.T) {
	repo := &mockRuleRepo{}
	engine := NewEngine(repo)

	rule1 := &model.Rule{
		ID:         "1",
		Enabled:    true,
		Type:       model.RuleMock,
		URLPattern: "/api/test",
		Method:     "GET",
	}
	engine.AddRule(rule1)

	tests := []struct {
		name   string
		method string
		url    string
		want   *model.Rule
	}{
		{"Exact match", "GET", "http://example.com/api/test", rule1},
		{"Method mismatch", "POST", "http://example.com/api/test", nil},
		{"URL mismatch", "GET", "http://example.com/api/other", nil},
		{"Partial pattern match", "GET", "http://example.com/v1/api/test/details", rule1},
		{"Any method", "PUT", "http://test.com/api/test", nil}, // method is set to GET in rule1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.url, nil)
			if got := engine.Match(req); got != tt.want {
				if got == nil || tt.want == nil || got.ID != tt.want.ID {
					t.Errorf("Engine.Match() = %v, want %v", got, tt.want)
				}
			}
		})
	}

	// Test case for disabled rule
	rule1.Enabled = false
	reqDisabled, _ := http.NewRequest("GET", "http://example.com/api/test", nil)
	if got := engine.Match(reqDisabled); got != nil {
		t.Errorf("Expected no match for disabled rule")
	}
	rule1.Enabled = true // Reset

	// Test case where method is empty
	rule2 := &model.Rule{ID: "2", URLPattern: "test", Enabled: true}
	repo.rules = []*model.Rule{rule2}
	reqEmptyMethod, _ := http.NewRequest("PATCH", "http://example.com/test", nil)
	if got := engine.Match(reqEmptyMethod); got == nil || got.ID != "2" {
		t.Errorf("Expected match for empty method")
	}

	// Test repo error
	repo.err = errors.New("repo error")
	if got := engine.Match(reqEmptyMethod); got != nil {
		t.Errorf("Expected nil on repo error")
	}
}

func TestEngine_UpdateAndDelete(t *testing.T) {
	repo := &mockRuleRepo{}
	engine := NewEngine(repo)

	rule := &model.Rule{ID: "1", URLPattern: "/old"}
	engine.AddRule(rule)

	rule.URLPattern = "/new"
	engine.UpdateRule(rule)
	if engine.GetRules()[0].URLPattern != "/new" {
		t.Errorf("Update failed")
	}

	engine.DeleteRule("1")
	if len(engine.GetRules()) != 0 {
		t.Errorf("Delete failed")
	}
}

func TestEngine_ErrorCases(t *testing.T) {
	repo := &mockRuleRepo{err: errors.New("db error")}
	engine := NewEngine(repo)

	// Suppress log output during test
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	engine.AddRule(&model.Rule{ID: "1"})
	engine.UpdateRule(&model.Rule{ID: "1"})
	engine.DeleteRule("1")

	if rules := engine.GetRules(); len(rules) != 0 {
		t.Errorf("Expected empty rules on error")
	}
}

func TestEngine_ClearRules(t *testing.T) {
	repo := &mockRuleRepo{}
	engine := NewEngine(repo)

	engine.AddRule(&model.Rule{ID: "1"})
	engine.AddRule(&model.Rule{ID: "2"})

	engine.ClearRules()
	if len(engine.GetRules()) != 0 {
		t.Errorf("Rules not cleared")
	}

	// Error case
	repo.err = errors.New("clear error")
	log.SetOutput(io.Discard)
	engine.ClearRules()
}

func TestEngine_MatchEdgeCases(t *testing.T) {
	repo := &mockRuleRepo{}
	engine := NewEngine(repo)

	t.Run("Match Any Method", func(t *testing.T) {
		rule := &model.Rule{ID: "any", URLPattern: "test", Enabled: true}
		repo.rules = []*model.Rule{rule}
		req, _ := http.NewRequest("POST", "http://test.com", nil)
		if got := engine.Match(req); got == nil || got.ID != "any" {
			t.Error("Expected match for any method")
		}
	})

	t.Run("Empty Pattern Match", func(t *testing.T) {
		rule := &model.Rule{ID: "empty", URLPattern: "", Enabled: true}
		repo.rules = []*model.Rule{rule}
		req, _ := http.NewRequest("GET", "http://any.com", nil)
		if got := engine.Match(req); got == nil || got.ID != "empty" {
			t.Error("Expected match for empty pattern")
		}
	})

	t.Run("No Rules", func(t *testing.T) {
		repo.rules = nil
		req, _ := http.NewRequest("GET", "http://any.com", nil)
		if got := engine.Match(req); got != nil {
			t.Error("Expected no match")
		}
	})
}
