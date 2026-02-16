package repository

import (
	"agent-proxy/internal/model"
	"database/sql"
	"encoding/json"
	"time"
)

type sqliteConfigRepository struct {
	db *sql.DB
}

func NewSQLiteConfigRepository(db *sql.DB) ConfigRepository {
	return &sqliteConfigRepository{db: db}
}

func (r *sqliteConfigRepository) Get() (*model.Config, error) {
	cfg := &model.Config{
		ProxyAddr:       ":8000",
		APIAddr:         ":8081",
		MCPAddr:         ":8082",
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
	db *sql.DB
}

func NewSQLiteTrafficRepository(db *sql.DB) TrafficRepository {
	return &sqliteTrafficRepository{db: db}
}

func (r *sqliteTrafficRepository) Add(entry *model.TrafficEntry) error {
	reqHeaders, _ := json.Marshal(entry.RequestHeaders)
	resHeaders, _ := json.Marshal(entry.ResponseHeaders)

	_, err := r.db.Exec(`
		INSERT INTO traffic (
			id, method, url, request_headers, request_body,
			status, response_headers, response_body, start_time, duration, modified_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entry.ID, entry.Method, entry.URL, string(reqHeaders), entry.RequestBody,
		entry.Status, string(resHeaders), entry.ResponseBody, entry.StartTime, int64(entry.Duration), entry.ModifiedBy)

	return err
}

func (r *sqliteTrafficRepository) GetPage(offset, limit int) ([]*model.TrafficEntry, int, error) {
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM traffic").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT 
			id, method, url, request_headers, request_body,
			status, response_headers, response_body, start_time, duration, modified_by
		FROM traffic ORDER BY start_time DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

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
		json.Unmarshal([]byte(reqH), &e.RequestHeaders)
		json.Unmarshal([]byte(resH), &e.ResponseHeaders)
		e.Duration = time.Duration(duration)
		entries = append(entries, &e)
	}
	return entries, total, nil
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

type sqliteRuleRepository struct {
	db *sql.DB
}

func NewSQLiteRuleRepository(db *sql.DB) RuleRepository {
	return &sqliteRuleRepository{db: db}
}

func (r *sqliteRuleRepository) GetAll() ([]*model.Rule, error) {
	rows, err := r.db.Query("SELECT id, type, url_pattern, method, strategy, response_json FROM rules")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*model.Rule
	for rows.Next() {
		var rule model.Rule
		var respJSON sql.NullString
		err := rows.Scan(&rule.ID, &rule.Type, &rule.URLPattern, &rule.Method, &rule.Strategy, &respJSON)
		if err != nil {
			continue
		}
		if respJSON.Valid && respJSON.String != "" {
			json.Unmarshal([]byte(respJSON.String), &rule.Response)
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
