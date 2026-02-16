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
		ProxyAddr:  ":8000",
		APIAddr:    ":8081",
		MCPAddr:    ":8082",
		MCPEnabled: false,
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
			status, response_headers, response_body, start_time, duration
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entry.ID, entry.Method, entry.URL, string(reqHeaders), entry.RequestBody,
		entry.Status, string(resHeaders), entry.ResponseBody, entry.StartTime, int64(entry.Duration))

	return err
}

func (r *sqliteTrafficRepository) GetAll() ([]*model.TrafficEntry, error) {
	rows, err := r.db.Query(`
		SELECT 
			id, method, url, request_headers, request_body,
			status, response_headers, response_body, start_time, duration
		FROM traffic ORDER BY start_time ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*model.TrafficEntry
	for rows.Next() {
		var e model.TrafficEntry
		var reqH, resH string
		var duration int64
		err := rows.Scan(
			&e.ID, &e.Method, &e.URL, &reqH, &e.RequestBody,
			&e.Status, &resH, &e.ResponseBody, &e.StartTime, &duration)
		if err != nil {
			continue
		}
		json.Unmarshal([]byte(reqH), &e.RequestHeaders)
		json.Unmarshal([]byte(resH), &e.ResponseHeaders)
		e.Duration = time.Duration(duration)
		entries = append(entries, &e)
	}
	return entries, nil
}

func (r *sqliteTrafficRepository) Clear() error {
	_, err := r.db.Exec("DELETE FROM traffic")
	return err
}
