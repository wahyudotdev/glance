package service

import (
	"fmt"
	"glance/internal/model"
	"testing"
)

type mockScenarioRepo struct {
	scenarios map[string]*model.Scenario
}

func (m *mockScenarioRepo) GetAll() ([]*model.Scenario, error) {
	var result []*model.Scenario
	for _, s := range m.scenarios {
		result = append(result, s)
	}
	return result, nil
}

func (m *mockScenarioRepo) GetByID(id string) (*model.Scenario, error) {
	return m.scenarios[id], nil
}

func (m *mockScenarioRepo) Add(s *model.Scenario) error {
	m.scenarios[s.ID] = s
	return nil
}

func (m *mockScenarioRepo) Update(s *model.Scenario) error {
	m.scenarios[s.ID] = s
	return nil
}

func (m *mockScenarioRepo) Delete(id string) error {
	delete(m.scenarios, id)
	return nil
}

type errorScenarioRepo struct {
	mockScenarioRepo
}

func (e *errorScenarioRepo) Add(_ *model.Scenario) error { return fmt.Errorf("error") }

func TestScenarioService_Errors(t *testing.T) {
	repo := &errorScenarioRepo{}
	svc := NewScenarioService(repo)

	err := svc.Create(&model.Scenario{Name: "Fail"})
	if err == nil {
		t.Error("Expected error on Create")
	}
}

func TestScenarioService(t *testing.T) {
	repo := &mockScenarioRepo{scenarios: make(map[string]*model.Scenario)}
	svc := NewScenarioService(repo)

	// Test Create
	s := &model.Scenario{Name: "Test"}
	err := svc.Create(s)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if s.ID == "" {
		t.Error("Expected ID to be generated")
	}

	// Test GetAll
	all, _ := svc.GetAll()
	if len(all) != 1 {
		t.Errorf("Expected 1 scenario, got %d", len(all))
	}

	// Test GetByID
	got, _ := svc.GetByID(s.ID)
	if got.Name != "Test" {
		t.Errorf("Expected name Test, got %s", got.Name)
	}

	// Test Update
	s.Name = "Updated"
	_ = svc.Update(s.ID, s)
	got, _ = svc.GetByID(s.ID)
	if got.Name != "Updated" {
		t.Errorf("Expected name Updated, got %s", got.Name)
	}

	// Test Delete
	_ = svc.Delete(s.ID)
	all, _ = svc.GetAll()
	if len(all) != 0 {
		t.Errorf("Expected 0 scenarios, got %d", len(all))
	}
}

func TestScenarioService_Mappings(t *testing.T) {
	repo := &mockScenarioRepo{scenarios: make(map[string]*model.Scenario)}
	svc := NewScenarioService(repo)

	s := &model.Scenario{
		Name: "With Mappings",
		VariableMappings: []model.VariableMapping{
			{Name: "token", SourceEntryID: "t1", SourcePath: "body.token", TargetJSONPath: "header.Auth"},
		},
	}
	_ = svc.Create(s)

	got, _ := svc.GetByID(s.ID)
	if len(got.VariableMappings) != 1 {
		t.Errorf("Expected 1 mapping")
	}
}

func TestScenarioService_Create_WithID(t *testing.T) {
	repo := &mockScenarioRepo{scenarios: make(map[string]*model.Scenario)}
	svc := NewScenarioService(repo)

	s := &model.Scenario{ID: "manual-id", Name: "Test"}
	_ = svc.Create(s)
	if s.ID != "manual-id" {
		t.Errorf("Expected manual-id, got %s", s.ID)
	}
}
