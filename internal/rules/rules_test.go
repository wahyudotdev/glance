package rules

import (
	"agent-proxy/internal/model"
	"net/http"
	"testing"
)

type mockRuleRepo struct {
	rules []*model.Rule
}

func (m *mockRuleRepo) GetAll() ([]*model.Rule, error) { return m.rules, nil }
func (m *mockRuleRepo) Add(r *model.Rule) error        { m.rules = append(m.rules, r); return nil }
func (m *mockRuleRepo) Update(r *model.Rule) error {
	for i, existing := range m.rules {
		if existing.ID == r.ID {
			m.rules[i] = r
			break
		}
	}
	return nil
}
func (m *mockRuleRepo) Delete(id string) error {
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
}
