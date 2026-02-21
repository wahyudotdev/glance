package apiserver

import (
	"glance/internal/mcp"
	"glance/internal/model"
	"glance/internal/service"
)

type mockConfigService struct {
	status map[string]any
	cfg    *model.Config
}

func (m *mockConfigService) GetStatus(_ *mcp.Server, _ string) (map[string]any, error) {
	return m.status, nil
}
func (m *mockConfigService) GetConfig() (*model.Config, error) { return m.cfg, nil }
func (m *mockConfigService) SaveConfig(cfg *model.Config) error {
	m.cfg = cfg
	return nil
}

type mockRuleService struct {
	rules []*model.Rule
}

func (m *mockRuleService) GetAll() []*model.Rule { return m.rules }
func (m *mockRuleService) Create(r *model.Rule) error {
	m.rules = append(m.rules, r)
	return nil
}
func (m *mockRuleService) Update(_ string, _ *model.Rule) error { return nil }
func (m *mockRuleService) Delete(_ string)                      {}

type mockTrafficService struct {
	entries []*model.TrafficEntry
}

func (m *mockTrafficService) GetPage(_, _ int) ([]*model.TrafficEntry, int) {
	return m.entries, len(m.entries)
}
func (m *mockTrafficService) Clear() {}

type mockScenarioService struct {
	scenarios []*model.Scenario
}

func (m *mockScenarioService) GetAll() ([]*model.Scenario, error) { return m.scenarios, nil }
func (m *mockScenarioService) GetByID(id string) (*model.Scenario, error) {
	for _, s := range m.scenarios {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, nil
}
func (m *mockScenarioService) Create(s *model.Scenario) error {
	m.scenarios = append(m.scenarios, s)
	return nil
}
func (m *mockScenarioService) Update(_ string, _ *model.Scenario) error { return nil }
func (m *mockScenarioService) Delete(_ string) error                    { return nil }

type mockInterceptService struct{}

func (m *mockInterceptService) ContinueRequest(_ string, _ service.ContinueRequestParams) error {
	return nil
}
func (m *mockInterceptService) ContinueResponse(_ string, _ service.ContinueResponseParams) error {
	return nil
}
func (m *mockInterceptService) Abort(_ string) error { return nil }
