package db

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	_ "modernc.org/sqlite"
)

func TestEnsureParentDir(t *testing.T) {
	if err := ensureParentDir("local.db"); err != nil {
		t.Fatalf("expected no error for local path: %v", err)
	}

	dirPath := filepath.Join(t.TempDir(), "a", "b", "c.db")
	if err := ensureParentDir(dirPath); err != nil {
		t.Fatalf("expected dir creation: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(dirPath)); err != nil {
		t.Fatalf("expected dir exists: %v", err)
	}

	if err := ensureParentDir("/proc/1/forbidden/db.sqlite"); err == nil {
		t.Fatalf("expected ensureParentDir permission error")
	}
}

func TestEnsureSchema(t *testing.T) {
	db, err := sql.Open("sqlite", "file:"+filepath.Join(t.TempDir(), "schema.db"))
	if err != nil {
		t.Fatalf("open sqlite error: %v", err)
	}
	if err := ensureSchema(context.Background(), db); err != nil {
		t.Fatalf("ensure schema error: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}
	if err := ensureSchema(context.Background(), db); err == nil {
		t.Fatalf("expected schema error on closed db")
	}
}

func TestNewSQLiteDB(t *testing.T) {
	db, err := NewSQLiteDB(context.Background(), filepath.Join(t.TempDir(), "x", "app.db"))
	if err != nil {
		t.Fatalf("NewSQLiteDB error: %v", err)
	}

	var count int
	if err := db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE type='table' AND name='todos'`).Scan(&count); err != nil {
		t.Fatalf("query schema error: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected todos table exists")
	}

	if _, err := NewSQLiteDB(context.Background(), "/proc/1/forbidden/db.sqlite"); err == nil {
		t.Fatalf("expected NewSQLiteDB error")
	}
}

func TestNewSQLiteDBOpenError(t *testing.T) {
	patches := gomonkey.ApplyFunc(sql.Open, func(string, string) (*sql.DB, error) {
		return nil, errors.New("open error")
	})
	defer patches.Reset()

	if _, err := NewSQLiteDB(context.Background(), filepath.Join(t.TempDir(), "x", "open-error.db")); err == nil {
		t.Fatalf("expected open error")
	}
}

func TestNewSQLiteDBPingError(t *testing.T) {
	if _, err := NewSQLiteDB(context.Background(), t.TempDir()); err == nil {
		t.Fatalf("expected ping error")
	}
}

func TestNewSQLiteDBEnsureSchemaError(t *testing.T) {
	if _, err := NewSQLiteDB(context.Background(), "/dev/null"); err == nil {
		t.Fatalf("expected ensure schema error")
	}
}
