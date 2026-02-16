package rules

import (
	"agent-proxy/internal/model"
	"agent-proxy/internal/repository"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Engine struct {
	mu   sync.RWMutex
	repo repository.RuleRepository
}

func NewEngine(repo repository.RuleRepository) *Engine {
	return &Engine{
		repo: repo,
	}
}

func (e *Engine) AddRule(rule *model.Rule) {
	if err := e.repo.Add(rule); err != nil {
		log.Printf("Error persisting rule: %v", err)
	}
}

func (e *Engine) GetRules() []*model.Rule {
	rules, err := e.repo.GetAll()
	if err != nil {
		log.Printf("Error loading rules: %v", err)
		return []*model.Rule{}
	}
	return rules
}

func (e *Engine) ClearRules() {
	// Not implemented in repo yet, but we can iterate or add a Clear method to RuleRepo
}

func (e *Engine) DeleteRule(id string) {
	if err := e.repo.Delete(id); err != nil {
		log.Printf("Error deleting rule: %v", err)
	}
}

func (e *Engine) UpdateRule(rule *model.Rule) {
	if err := e.repo.Update(rule); err != nil {
		log.Printf("Error updating rule: %v", err)
	}
}

func (e *Engine) Match(r *http.Request) *model.Rule {
	e.mu.RLock()
	// We load from repo every time for now to keep it simple and consistent,
	// but we could cache them in memory.
	rules, err := e.repo.GetAll()
	e.mu.RUnlock()

	if err != nil {
		return nil
	}

	for _, rule := range rules {
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
