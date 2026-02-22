// Package rules implements the matching logic for mocks and breakpoints.
package rules

import (
	"glance/internal/model"
	"glance/internal/repository"
	"log"
	"net/http"
	"strings"
	"sync"
)

// Engine manages the collection of active interception rules.
type Engine struct {
	mu   sync.RWMutex
	repo repository.RuleRepository
}

// NewEngine creates a new Engine with the provided rule repository.
func NewEngine(repo repository.RuleRepository) *Engine {
	return &Engine{
		repo: repo,
	}
}

// AddRule adds a new rule to the engine and persists it.
func (e *Engine) AddRule(rule *model.Rule) {
	if err := e.repo.Add(rule); err != nil {
		log.Printf("Error persisting rule: %v", err)
	}
}

// GetRules retrieves all active rules from the repository.
func (e *Engine) GetRules() []*model.Rule {
	rules, err := e.repo.GetAll()
	if err != nil {
		log.Printf("Error loading rules: %v", err)
		return []*model.Rule{}
	}
	return rules
}

// ClearRules removes all active rules from the repository.
func (e *Engine) ClearRules() {
	rules, err := e.repo.GetAll()
	if err != nil {
		log.Printf("Error loading rules for clearing: %v", err)
		return
	}
	for _, r := range rules {
		if err := e.repo.Delete(r.ID); err != nil {
			log.Printf("Error deleting rule %s: %v", r.ID, err)
		}
	}
}

// DeleteRule removes a rule by its ID and updates the repository.
func (e *Engine) DeleteRule(id string) {
	if err := e.repo.Delete(id); err != nil {
		log.Printf("Error deleting rule: %v", err)
	}
}

// UpdateRule modifies an existing rule and persists the changes.
func (e *Engine) UpdateRule(rule *model.Rule) {
	if err := e.repo.Update(rule); err != nil {
		log.Printf("Error updating rule: %v", err)
	}
}

// Match checks if an incoming HTTP request matches any active rules.
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
		if !rule.Enabled {
			continue
		}
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
