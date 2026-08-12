package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const businessSchemaVersion = 1

var errRecordNotFound = errors.New("business record not found")
var errRecordConflict = errors.New("business record already exists")

type businessStore struct {
	db   *sql.DB
	path string
}

type storedRecord struct {
	Kind      string
	ID        string
	Data      json.RawMessage
	CreatedAt string
	UpdatedAt string
}

func newBusinessStore() (*businessStore, error) {
	path := strings.TrimSpace(os.Getenv("STA100_DB_PATH"))
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, ".local", "share", "sta100", "sta100.db")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	dsn := "file:" + filepath.ToSlash(path) + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	store := &businessStore{db: db, path: path}
	if err := store.migrate(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	if err := os.Chmod(path, 0600); err != nil {
		db.Close()
		return nil, err
	}
	if err := store.seed(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (s *businessStore) Close() error { return s.db.Close() }

func (s *businessStore) migrate(ctx context.Context) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS schema_meta (version INTEGER NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS records (
			kind TEXT NOT NULL,
			id TEXT NOT NULL,
			data_json TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			deleted_at TEXT,
			PRIMARY KEY (kind, id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_records_kind_updated ON records(kind, updated_at DESC)`,
		`CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value_json TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			action TEXT NOT NULL,
			entity_type TEXT NOT NULL,
			entity_id TEXT NOT NULL,
			result TEXT NOT NULL,
			operator TEXT NOT NULL,
			details_json TEXT NOT NULL,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS agent_token_usage (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			request_id TEXT NOT NULL,
			page TEXT NOT NULL,
			feature TEXT NOT NULL,
			stage TEXT NOT NULL,
			agent_id TEXT NOT NULL,
			run_id TEXT NOT NULL,
			status TEXT NOT NULL,
			input_tokens INTEGER NOT NULL DEFAULT 0,
			output_tokens INTEGER NOT NULL DEFAULT 0,
			cache_read_tokens INTEGER NOT NULL DEFAULT 0,
			cache_write_tokens INTEGER NOT NULL DEFAULT 0,
			total_tokens INTEGER NOT NULL DEFAULT 0,
			usage_available INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_agent_token_usage_created ON agent_token_usage(created_at DESC)`,
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM schema_meta`).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_meta(version) VALUES (?)`, businessSchemaVersion); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *businessStore) create(ctx context.Context, kind, id string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err = s.db.ExecContext(ctx, `INSERT INTO records(kind,id,data_json,created_at,updated_at) VALUES(?,?,?,?,?)`, kind, id, data, now, now)
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "unique") {
		return errRecordConflict
	}
	return err
}

func (s *businessStore) put(ctx context.Context, kind, id string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	result, err := s.db.ExecContext(ctx, `UPDATE records SET data_json=?,updated_at=?,deleted_at=NULL WHERE kind=? AND id=?`, data, now, kind, id)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed == 0 {
		return errRecordNotFound
	}
	return nil
}

func (s *businessStore) get(ctx context.Context, kind, id string, destination any) error {
	var data []byte
	err := s.db.QueryRowContext(ctx, `SELECT data_json FROM records WHERE kind=? AND id=? AND deleted_at IS NULL`, kind, id).Scan(&data)
	if errors.Is(err, sql.ErrNoRows) {
		return errRecordNotFound
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, destination)
}

func listRecords[T any](ctx context.Context, s *businessStore, kind string) ([]T, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT data_json FROM records WHERE kind=? AND deleted_at IS NULL ORDER BY updated_at DESC`, kind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]T, 0)
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}
		var item T
		if err := json.Unmarshal(data, &item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *businessStore) softDelete(ctx context.Context, kind, id string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	result, err := s.db.ExecContext(ctx, `UPDATE records SET deleted_at=?,updated_at=? WHERE kind=? AND id=? AND deleted_at IS NULL`, now, now, kind, id)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed == 0 {
		return errRecordNotFound
	}
	return nil
}

func (s *businessStore) nextSequence(ctx context.Context, kind, prefix string, width int) (string, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM records WHERE kind=?`, kind).Scan(&count); err != nil {
		return "", err
	}
	for offset := 1; offset < 100000; offset++ {
		candidate := fmt.Sprintf("%s-%0*d", prefix, width, count+offset)
		var exists int
		if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM records WHERE kind=? AND id=?`, kind, candidate).Scan(&exists); err != nil {
			return "", err
		}
		if exists == 0 {
			return candidate, nil
		}
	}
	return "", errors.New("sequence exhausted")
}

func (s *businessStore) getSetting(ctx context.Context, key string, destination any) error {
	var data []byte
	err := s.db.QueryRowContext(ctx, `SELECT value_json FROM settings WHERE key=?`, key).Scan(&data)
	if errors.Is(err, sql.ErrNoRows) {
		return errRecordNotFound
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, destination)
}

func (s *businessStore) putSetting(ctx context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err = s.db.ExecContext(ctx, `INSERT INTO settings(key,value_json,updated_at) VALUES(?,?,?)
		ON CONFLICT(key) DO UPDATE SET value_json=excluded.value_json,updated_at=excluded.updated_at`, key, data, now)
	return err
}

func (s *businessStore) audit(ctx context.Context, action, entityType, entityID, operator string, details any) {
	data, err := json.Marshal(details)
	if err != nil {
		data = []byte(`{}`)
	}
	_, _ = s.db.ExecContext(ctx, `INSERT INTO audit_logs(action,entity_type,entity_id,result,operator,details_json,created_at) VALUES(?,?,?,?,?,?,?)`,
		action, entityType, entityID, "success", operator, data, time.Now().UTC().Format(time.RFC3339Nano))
}

func timeNowUTC() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}
