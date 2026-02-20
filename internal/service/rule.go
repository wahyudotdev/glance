// Package service implements the core business logic.
package service

import (
	"glance/internal/model"
	"glance/internal/rules"

	"github.com/google/uuid"
)

// RuleService defines the interface for managing interception rules.
type RuleService interface {
	GetAll() []*model.Rule
	Create(rule *model.Rule) error
	Update(id string, rule *model.Rule) error
	Delete(id string)
}

type ruleService struct {
	engine *rules.Engine
}

// NewRuleService creates a new RuleService.
func NewRuleService(engine *rules.Engine) RuleService {
	return &ruleService{engine: engine}
}

func (s *ruleService) GetAll() []*model.Rule {
	return s.engine.GetRules()
}

func (s *ruleService) Create(rule *model.Rule) error {
	if rule.ID == "" {
		rule.ID = uuid.New().String()
	}
	s.engine.AddRule(rule)
	return nil
}

func (s *ruleService) Update(id string, rule *model.Rule) error {
	rule.ID = id
	s.engine.UpdateRule(rule)
	return nil
}

func (s *ruleService) Delete(id string) {
	s.engine.DeleteRule(id)
}
