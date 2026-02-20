package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"glance/internal/model"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type sqliteConfigRepository struct {
	db *sql.DB
}

// NewSQLiteConfigRepository creates a new SQLite-backed ConfigRepository.
func NewSQLiteConfigRepository(db *sql.DB) ConfigRepository {
	return &sqliteConfigRepository{db: db}
}

func (r *sqliteConfigRepository) Get() (*model.Config, error) {
	cfg := &model.Config{
		ProxyAddr:       ":15500",
		APIAddr:         ":15501",
		MCPAddr:         ":15502",
		MCPEnabled:      false,
		HistoryLimit:    500,
		MaxResponseSize: 1024 * 1024, // 1 MB
		DefaultPageSize: 50,
	}

	var val string
	err := r.db.QueryRow("SELECT value FROM config WHERE key = 'app_config'").Scan(&val)
	if err == nil {
		if err := json.Unmarshal([]byte(val), cfg); err != nil {
			return nil, err
		}
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	return cfg, nil
}

func (r *sqliteConfigRepository) Save(cfg *model.Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("INSERT OR REPLACE INTO config (key, value) VALUES ('app_config', ?)", string(data))
	return err
}

type sqliteTrafficRepository struct {
	db         *sql.DB
	writeQueue chan *model.TrafficEntry
	memCache   []*model.TrafficEntry
	cacheSize  int
	mu         sync.RWMutex
}

// NewSQLiteTrafficRepository creates a new SQLite-backed TrafficRepository.
func NewSQLiteTrafficRepository(db *sql.DB) TrafficRepository {
	repo := &sqliteTrafficRepository{
		db:         db,
		writeQueue: make(chan *model.TrafficEntry, 100),
		memCache:   make([]*model.TrafficEntry, 0, 500),
		cacheSize:  500,
	}
	go repo.writeWorker()
	return repo
}

func (r *sqliteTrafficRepository) writeWorker() {
	for entry := range r.writeQueue {
		reqHeaders, _ := json.Marshal(entry.RequestHeaders)
		resHeaders, _ := json.Marshal(entry.ResponseHeaders)

		_, err := r.db.Exec(`
			INSERT INTO traffic (
				id, method, url, request_headers, request_body,
				status, response_headers, response_body, start_time, duration, modified_by
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			entry.ID, entry.Method, entry.URL, string(reqHeaders), entry.RequestBody,
			entry.Status, string(resHeaders), entry.ResponseBody, entry.StartTime, int64(entry.Duration), entry.ModifiedBy)

		if err != nil {
			log.Printf("Background DB write error: %v", err)
		}
	}
}

func (r *sqliteTrafficRepository) Add(entry *model.TrafficEntry) error {
	// 1. Update Memory Cache immediately for fast UI response
	r.mu.Lock()
	r.memCache = append(r.memCache, entry)
	if len(r.memCache) > r.cacheSize {
		r.memCache = r.memCache[1:]
	}
	r.mu.Unlock()

	// 2. Queue for background persistent storage
	select {
	case r.writeQueue <- entry:
	default:
		log.Printf("Warning: Traffic write queue full, dropping entry %s", entry.ID)
	}
	return nil
}

func (r *sqliteTrafficRepository) GetPage(offset, limit int) ([]*model.TrafficEntry, int, error) {
	// If we are asking for the first page and it might be in memory, we can optimize.
	// But for now, to keep consistency with total count, we still query DB total.
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM traffic").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// If offset is 0, we can prepend/merge with memory cache for "un-flushed" items
	// To keep it simple and fix the 'locked' error, we mainly needed serialized writes.
	// Let's stick to DB for GetPage but the 'locked' error will be gone because
	// writes are now serialized in the background worker.

	rows, err := r.db.Query(`
		SELECT 
			id, method, url, request_headers, request_body,
			status, response_headers, response_body, start_time, duration, modified_by
		FROM traffic ORDER BY start_time DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var entries []*model.TrafficEntry
	for rows.Next() {
		var e model.TrafficEntry
		var reqH, resH string
		var duration int64
		err := rows.Scan(
			&e.ID, &e.Method, &e.URL, &reqH, &e.RequestBody,
			&e.Status, &resH, &e.ResponseBody, &e.StartTime, &duration, &e.ModifiedBy)
		if err != nil {
			continue
		}
		_ = json.Unmarshal([]byte(reqH), &e.RequestHeaders)
		_ = json.Unmarshal([]byte(resH), &e.ResponseHeaders)
		e.Duration = time.Duration(duration)
		entries = append(entries, &e)
	}
	return entries, total, nil
}

func (r *sqliteTrafficRepository) GetByIDs(ids []string) ([]*model.TrafficEntry, error) {
	if len(ids) == 0 {
		return []*model.TrafficEntry{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT 
			id, method, url, request_headers, request_body,
			status, response_headers, response_body, start_time, duration, modified_by
		FROM traffic WHERE id IN (%s)`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var entries []*model.TrafficEntry
	for rows.Next() {
		var e model.TrafficEntry
		var reqH, resH string
		var duration int64
		err := rows.Scan(
			&e.ID, &e.Method, &e.URL, &reqH, &e.RequestBody,
			&e.Status, &resH, &e.ResponseBody, &e.StartTime, &duration, &e.ModifiedBy)
		if err != nil {
			continue
		}
		_ = json.Unmarshal([]byte(reqH), &e.RequestHeaders)
		_ = json.Unmarshal([]byte(resH), &e.ResponseHeaders)
		e.Duration = time.Duration(duration)
		entries = append(entries, &e)
	}
	return entries, nil
}

func (r *sqliteTrafficRepository) Clear() error {
	_, err := r.db.Exec("DELETE FROM traffic")
	return err
}

func (r *sqliteTrafficRepository) Prune(limit int) error {
	_, err := r.db.Exec(`
		DELETE FROM traffic WHERE id NOT IN (
			SELECT id FROM traffic ORDER BY start_time DESC LIMIT ?
		)`, limit)
	return err
}

func (r *sqliteTrafficRepository) Flush() {
	// Simple flush: send a "no-op" entry and wait for it if possible,
	// or just sleep briefly. A better way is using a specialized signal.
	// For now, since it's a test helper, we'll use a small sleep or
	// we can improve the worker to handle a "sync" command.
	time.Sleep(100 * time.Millisecond)
}

type sqliteRuleRepository struct {
	db *sql.DB
}

// NewSQLiteRuleRepository creates a new SQLite-backed RuleRepository.
func NewSQLiteRuleRepository(db *sql.DB) RuleRepository {
	return &sqliteRuleRepository{db: db}
}

func (r *sqliteRuleRepository) GetAll() ([]*model.Rule, error) {
	rows, err := r.db.Query("SELECT id, type, url_pattern, method, strategy, response_json FROM rules")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var rules []*model.Rule
	for rows.Next() {
		var rule model.Rule
		var respJSON sql.NullString
		err := rows.Scan(&rule.ID, &rule.Type, &rule.URLPattern, &rule.Method, &rule.Strategy, &respJSON)
		if err != nil {
			continue
		}
		if respJSON.Valid && respJSON.String != "" {
			_ = json.Unmarshal([]byte(respJSON.String), &rule.Response)
		}
		rules = append(rules, &rule)
	}
	return rules, nil
}

func (r *sqliteRuleRepository) Add(rule *model.Rule) error {
	respJSON, _ := json.Marshal(rule.Response)
	_, err := r.db.Exec(`
		INSERT INTO rules (id, type, url_pattern, method, strategy, response_json)
		VALUES (?, ?, ?, ?, ?, ?)`,
		rule.ID, rule.Type, rule.URLPattern, rule.Method, rule.Strategy, string(respJSON))
	return err
}

func (r *sqliteRuleRepository) Update(rule *model.Rule) error {
	respJSON, _ := json.Marshal(rule.Response)
	_, err := r.db.Exec(`
		UPDATE rules SET type = ?, url_pattern = ?, method = ?, strategy = ?, response_json = ?
		WHERE id = ?`,
		rule.Type, rule.URLPattern, rule.Method, rule.Strategy, string(respJSON), rule.ID)
	return err
}

func (r *sqliteRuleRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM rules WHERE id = ?", id)
	return err
}

type sqliteScenarioRepository struct {
	db *sql.DB
}

// NewSQLiteScenarioRepository creates a new SQLite-backed ScenarioRepository.
func NewSQLiteScenarioRepository(db *sql.DB) ScenarioRepository {
	return &sqliteScenarioRepository{db: db}
}

func (r *sqliteScenarioRepository) GetAll() ([]*model.Scenario, error) {
	rows, err := r.db.Query("SELECT id FROM scenarios ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	_ = rows.Close() // Close immediately to release the connection

	var scenarios []*model.Scenario
	for _, id := range ids {
		s, err := r.GetByID(id)
		if err != nil {
			log.Printf("Error loading scenario %s: %v", id, err)
			continue
		}
		scenarios = append(scenarios, s)
	}
	return scenarios, nil
}

func (r *sqliteScenarioRepository) GetByID(id string) (*model.Scenario, error) {
	s := &model.Scenario{}
	err := r.db.QueryRow("SELECT id, name, description, created_at FROM scenarios WHERE id = ?", id).
		Scan(&s.ID, &s.Name, &s.Description, &s.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Load steps with joined traffic data
	stepRows, err := r.db.Query(`
		SELECT 
			s.id, s.traffic_entry_id, s.step_order, s.notes,
			t.method, t.url, t.status
		FROM scenario_steps s
		LEFT JOIN traffic t ON s.traffic_entry_id = t.id
		WHERE s.scenario_id = ? 
		ORDER BY s.step_order ASC`, id)
	if err == nil {
		defer func() { _ = stepRows.Close() }()
		for stepRows.Next() {
			var step model.ScenarioStep
			var method, url sql.NullString
			var status sql.NullInt32
			err := stepRows.Scan(&step.ID, &step.TrafficEntryID, &step.Order, &step.Notes, &method, &url, &status)
			if err == nil {
				if method.Valid {
					step.TrafficEntry = &model.TrafficEntry{
						ID:     step.TrafficEntryID,
						Method: method.String,
						URL:    url.String,
						Status: int(status.Int32),
					}
				}
				s.Steps = append(s.Steps, step)
			}
		}
	}

	// Load mappings
	mappingRows, err := r.db.Query("SELECT name, source_entry_id, source_path, target_json_path FROM variable_mappings WHERE scenario_id = ?", id)
	if err == nil {
		defer func() { _ = mappingRows.Close() }()
		for mappingRows.Next() {
			var m model.VariableMapping
			if err := mappingRows.Scan(&m.Name, &m.SourceEntryID, &m.SourcePath, &m.TargetJSONPath); err == nil {
				s.VariableMappings = append(s.VariableMappings, m)
			}
		}
	}

	return s, nil
}

func (r *sqliteScenarioRepository) Add(s *model.Scenario) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.Exec("INSERT INTO scenarios (id, name, description, created_at) VALUES (?, ?, ?, ?)",
		s.ID, s.Name, s.Description, s.CreatedAt)
	if err != nil {
		log.Printf("Error inserting scenario: %v", err)
		return err
	}

	for _, step := range s.Steps {
		if step.ID == "" {
			step.ID = uuid.New().String()
		}
		_, err = tx.Exec("INSERT INTO scenario_steps (id, scenario_id, traffic_entry_id, step_order, notes) VALUES (?, ?, ?, ?, ?)",
			step.ID, s.ID, step.TrafficEntryID, step.Order, step.Notes)
		if err != nil {
			log.Printf("Error inserting scenario step: %v", err)
			return err
		}
	}

	for _, m := range s.VariableMappings {
		_, err = tx.Exec("INSERT INTO variable_mappings (id, scenario_id, name, source_entry_id, source_path, target_json_path) VALUES (?, ?, ?, ?, ?, ?)",
			uuid.New().String(), s.ID, m.Name, m.SourceEntryID, m.SourcePath, m.TargetJSONPath)
		if err != nil {
			log.Printf("Error inserting variable mapping: %v", err)
			return err
		}
	}

	return tx.Commit()
}

func (r *sqliteScenarioRepository) Update(s *model.Scenario) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.Exec("UPDATE scenarios SET name = ?, description = ? WHERE id = ?", s.Name, s.Description, s.ID)
	if err != nil {
		log.Printf("Error updating scenario: %v", err)
		return err
	}

	// Simple approach: delete and re-insert steps and mappings
	_, _ = tx.Exec("DELETE FROM scenario_steps WHERE scenario_id = ?", s.ID)
	_, _ = tx.Exec("DELETE FROM variable_mappings WHERE scenario_id = ?", s.ID)

	for _, step := range s.Steps {
		if step.ID == "" {
			step.ID = uuid.New().String()
		}
		_, err = tx.Exec("INSERT INTO scenario_steps (id, scenario_id, traffic_entry_id, step_order, notes) VALUES (?, ?, ?, ?, ?)",
			step.ID, s.ID, step.TrafficEntryID, step.Order, step.Notes)
		if err != nil {
			log.Printf("Error re-inserting scenario step: %v", err)
			return err
		}
	}

	for _, m := range s.VariableMappings {
		_, err = tx.Exec("INSERT INTO variable_mappings (id, scenario_id, name, source_entry_id, source_path, target_json_path) VALUES (?, ?, ?, ?, ?, ?)",
			uuid.New().String(), s.ID, m.Name, m.SourceEntryID, m.SourcePath, m.TargetJSONPath)
		if err != nil {
			log.Printf("Error re-inserting variable mapping: %v", err)
			return err
		}
	}

	return tx.Commit()
}

func (r *sqliteScenarioRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM scenarios WHERE id = ?", id)
	return err
}
