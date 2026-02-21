package service

import (
	"glance/internal/model"
)

type mockTrafficRepo struct {
	entries []*model.TrafficEntry
}

func (m *mockTrafficRepo) Add(e *model.TrafficEntry) error {
	m.entries = append(m.entries, e)
	return nil
}

func (m *mockTrafficRepo) GetPage(offset, limit int) ([]*model.TrafficEntry, int, error) {
	start := offset
	if start > len(m.entries) {
		start = len(m.entries)
	}
	end := offset + limit
	if end > len(m.entries) {
		end = len(m.entries)
	}
	return m.entries[start:end], len(m.entries), nil
}

func (m *mockTrafficRepo) GetByIDs(_ []string) ([]*model.TrafficEntry, error) {
	return nil, nil
}

func (m *mockTrafficRepo) Clear() error {
	m.entries = nil
	return nil
}

func (m *mockTrafficRepo) Prune(_ int) error {
	return nil
}

func (m *mockTrafficRepo) Flush() {}

type mockRuleRepo struct {
	rules map[string]*model.Rule
}

func (m *mockRuleRepo) GetAll() ([]*model.Rule, error) {
	var res []*model.Rule
	for _, r := range m.rules {
		res = append(res, r)
	}
	return res, nil
}

func (m *mockRuleRepo) Add(r *model.Rule) error {
	m.rules[r.ID] = r
	return nil
}

func (m *mockRuleRepo) Update(r *model.Rule) error {
	m.rules[r.ID] = r
	return nil
}

func (m *mockRuleRepo) Delete(id string) error {
	delete(m.rules, id)
	return nil
}
