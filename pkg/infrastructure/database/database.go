package database

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type TodoRecord struct {
	ID        string
	Title     string
	Completed bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DB struct {
	sqlDB *sql.DB
}

var openSQL = sql.Open
var getRowsAffected = func(result sql.Result) (int64, error) {
	return result.RowsAffected()
}

func Connect(path string) *DB {
	db, err := Open(path)
	if err != nil {
		panic(err)
	}
	return db
}

func Open(path string) (*DB, error) {
	if err := ensureParentDir(path); err != nil {
		return nil, err
	}
	sqlDB, err := openSQL("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := ensureSchema(context.Background(), sqlDB); err != nil {
		return nil, errors.Join(err, sqlDB.Close())
	}
	return &DB{sqlDB: sqlDB}, nil
}

func (db *DB) Close() error {
	return db.sqlDB.Close()
}

func (db *DB) CreateTodo(ctx context.Context, record TodoRecord) (TodoRecord, error) {
	_, err := db.sqlDB.ExecContext(
		ctx,
		`INSERT INTO todos (id, title, completed, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		record.ID,
		record.Title,
		boolToInt(record.Completed),
		record.CreatedAt.Format(time.RFC3339Nano),
		record.UpdatedAt.Format(time.RFC3339Nano),
	)
	return record, err
}

func (db *DB) GetTodo(ctx context.Context, id string) (TodoRecord, bool, error) {
	row := db.sqlDB.QueryRowContext(ctx, `SELECT id, title, completed, created_at, updated_at FROM todos WHERE id = ?`, id)
	record, err := scanTodo(row)
	if errors.Is(err, sql.ErrNoRows) {
		return TodoRecord{}, false, nil
	}
	return record, err == nil, err
}

func (db *DB) ListTodos(ctx context.Context) ([]TodoRecord, error) {
	rows, err := db.sqlDB.QueryContext(ctx, `SELECT id, title, completed, created_at, updated_at FROM todos ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var records []TodoRecord
	for rows.Next() {
		record, err := scanTodo(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, rows.Err()
}

func (db *DB) UpdateTodo(ctx context.Context, record TodoRecord) (TodoRecord, bool, error) {
	result, err := db.sqlDB.ExecContext(
		ctx,
		`UPDATE todos SET title = ?, completed = ?, updated_at = ? WHERE id = ?`,
		record.Title,
		boolToInt(record.Completed),
		record.UpdatedAt.Format(time.RFC3339Nano),
		record.ID,
	)
	if err != nil {
		return TodoRecord{}, false, err
	}
	affected, err := getRowsAffected(result)
	if err != nil {
		return TodoRecord{}, false, err
	}
	return record, affected > 0, nil
}

func (db *DB) DeleteTodo(ctx context.Context, id string) (bool, error) {
	result, err := db.sqlDB.ExecContext(ctx, `DELETE FROM todos WHERE id = ?`, id)
	if err != nil {
		return false, err
	}
	affected, err := getRowsAffected(result)
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func ensureParentDir(path string) error {
	dir := filepath.Dir(path)
	if path == ":memory:" || dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

func ensureSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS todos (
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	completed INTEGER NOT NULL DEFAULT 0,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);`)
	return err
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTodo(s scanner) (TodoRecord, error) {
	var record TodoRecord
	var completed int
	var createdAt string
	var updatedAt string
	if err := s.Scan(&record.ID, &record.Title, &completed, &createdAt, &updatedAt); err != nil {
		return TodoRecord{}, err
	}
	parsedCreatedAt, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return TodoRecord{}, err
	}
	parsedUpdatedAt, err := time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return TodoRecord{}, err
	}
	record.Completed = completed != 0
	record.CreatedAt = parsedCreatedAt
	record.UpdatedAt = parsedUpdatedAt
	return record, nil
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
