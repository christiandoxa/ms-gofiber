package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

var (
	openSQLiteDB = func(driverName, dataSourceName string) (*sql.DB, error) {
		return sql.Open(driverName, dataSourceName)
	}
	pingSQLiteDB = func(ctx context.Context, db *sql.DB) error {
		return db.PingContext(ctx)
	}
	ensureSQLiteSchema = ensureSchema
)

type SQLiteOptions struct {
	Path string
}

func NewSQLiteDB(ctx context.Context, opts SQLiteOptions) (*sql.DB, error) {
	if err := ensureParentDir(opts.Path); err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)", opts.Path)
	db, err := openSQLiteDB("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pingSQLiteDB(pctx, db); err != nil {
		return nil, errors.Join(err, db.Close())
	}

	if err := ensureSQLiteSchema(ctx, db); err != nil {
		return nil, errors.Join(err, db.Close())
	}

	return db, nil
}

func ensureParentDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

func ensureSchema(ctx context.Context, db *sql.DB) error {
	const schema = `
CREATE TABLE IF NOT EXISTS todos (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    completed INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos (created_at DESC);
`
	_, err := db.ExecContext(ctx, schema)
	return err
}
