package apiserver

import (
	"glance/internal/mcp"
	"glance/internal/model"
	"glance/internal/service"
)

type mockConfigService struct {
	status map[string]any
	cfg    *model.Config
	err    error
}

func (m *mockConfigService) GetStatus(_ *mcp.Server, _ string) (map[string]any, error) {
	return m.status, m.err
}
func (m *mockConfigService) GetConfig() (*model.Config, error) {
	return m.cfg, m.err
}
func (m *mockConfigService) SaveConfig(cfg *model.Config) error {
	if m.err != nil {
		return m.err
	}
	m.cfg = cfg
	return nil
}

type mockRuleService struct {
	rules []*model.Rule
	err   error
}

func (m *mockRuleService) GetAll() []*model.Rule { return m.rules }
func (m *mockRuleService) Create(r *model.Rule) error {
	if m.err != nil {
		return m.err
	}
	m.rules = append(m.rules, r)
	return nil
}
func (m *mockRuleService) Update(_ string, r *model.Rule) error {
	if m.err != nil {
		return m.err
	}
	return nil
}
func (m *mockRuleService) Delete(_ string) {}

type mockTrafficService struct {
	entries []*model.TrafficEntry
}

func (m *mockTrafficService) GetPage(_, _ int) ([]*model.TrafficEntry, int) {
	return m.entries, len(m.entries)
}
func (m *mockTrafficService) Clear() {}

type mockScenarioService struct {
	scenarios []*model.Scenario
	err       error
}

func (m *mockScenarioService) GetAll() ([]*model.Scenario, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.scenarios, nil
}
func (m *mockScenarioService) GetByID(id string) (*model.Scenario, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, s := range m.scenarios {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, nil
}
func (m *mockScenarioService) Create(s *model.Scenario) error {
	if m.err != nil {
		return m.err
	}
	m.scenarios = append(m.scenarios, s)
	return nil
}
func (m *mockScenarioService) Update(_ string, _ *model.Scenario) error {
	return m.err
}
func (m *mockScenarioService) Delete(_ string) error {
	return m.err
}

type mockInterceptService struct {
	err error
}

func (m *mockInterceptService) ContinueRequest(_ string, _ service.ContinueRequestParams) error {
	return m.err
}
func (m *mockInterceptService) ContinueResponse(_ string, _ service.ContinueResponseParams) error {
	return m.err
}
func (m *mockInterceptService) Abort(_ string) error { return m.err }
