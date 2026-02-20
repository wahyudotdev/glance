package service

import (
	"glance/internal/model"
	"glance/internal/repository"
	"time"

	"github.com/google/uuid"
)

// ScenarioService defines the interface for managing scenarios.
type ScenarioService interface {
	GetAll() ([]*model.Scenario, error)
	GetByID(id string) (*model.Scenario, error)
	Create(s *model.Scenario) error
	Update(id string, s *model.Scenario) error
	Delete(id string) error
}

type scenarioService struct {
	repo repository.ScenarioRepository
}

// NewScenarioService creates a new ScenarioService.
func NewScenarioService(repo repository.ScenarioRepository) ScenarioService {
	return &scenarioService{repo: repo}
}

func (s *scenarioService) GetAll() ([]*model.Scenario, error) {
	return s.repo.GetAll()
}

func (s *scenarioService) GetByID(id string) (*model.Scenario, error) {
	return s.repo.GetByID(id)
}

func (s *scenarioService) Create(scenario *model.Scenario) error {
	if scenario.ID == "" {
		scenario.ID = uuid.New().String()
	}
	if scenario.CreatedAt.IsZero() {
		scenario.CreatedAt = time.Now()
	}
	return s.repo.Add(scenario)
}

func (s *scenarioService) Update(id string, scenario *model.Scenario) error {
	scenario.ID = id
	return s.repo.Update(scenario)
}

func (s *scenarioService) Delete(id string) error {
	return s.repo.Delete(id)
}
