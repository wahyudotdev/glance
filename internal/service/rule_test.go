package service

import (
	"glance/internal/model"
	"glance/internal/rules"
	"testing"
)

func TestRuleService(t *testing.T) {
	repo := &mockRuleRepo{rules: make(map[string]*model.Rule)}
	engine := rules.NewEngine(repo)
	svc := NewRuleService(engine)

	// Test Create
	rule := &model.Rule{
		Type:       model.RuleMock,
		URLPattern: "/test",
		Method:     "GET",
	}
	err := svc.Create(rule)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if rule.ID == "" {
		t.Error("Expected ID to be generated")
	}

	// Test GetAll
	all := svc.GetAll()
	if len(all) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(all))
	}

	// Test Update
	rule.URLPattern = "/updated"
	_ = svc.Update(rule.ID, rule)
	all = svc.GetAll()
	// Since it's a map, order is not guaranteed but we only have one
	if all[0].URLPattern != "/updated" {
		t.Errorf("Expected pattern /updated, got %s", all[0].URLPattern)
	}

	// Test Delete
	svc.Delete(rule.ID)
	all = svc.GetAll()
	if len(all) != 0 {
		t.Errorf("Expected 0 rules, got %d", len(all))
	}
}
