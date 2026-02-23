package db

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"ms-gofiber/internal/config"
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
	_ = db.Close()
	if err := ensureSchema(context.Background(), db); err == nil {
		t.Fatalf("expected schema error on closed db")
	}
}

func TestNewSQLiteDB(t *testing.T) {
	cfg := &config.Config{SQLitePath: filepath.Join(t.TempDir(), "x", "app.db")}
	db, err := NewSQLiteDB(context.Background(), cfg)
	if err != nil {
		t.Fatalf("NewSQLiteDB error: %v", err)
	}
	defer db.Close()

	var count int
	if err := db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE type='table' AND name='todos'`).Scan(&count); err != nil {
		t.Fatalf("query schema error: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected todos table exists")
	}

	bad := &config.Config{SQLitePath: "/proc/1/forbidden/db.sqlite"}
	if _, err := NewSQLiteDB(context.Background(), bad); err == nil {
		t.Fatalf("expected NewSQLiteDB error")
	}
}

func TestNewSQLiteDBOpenError(t *testing.T) {
	origOpen := openSQLiteDB
	t.Cleanup(func() { openSQLiteDB = origOpen })

	openSQLiteDB = func(string, string) (*sql.DB, error) {
		return nil, errors.New("open error")
	}

	cfg := &config.Config{SQLitePath: filepath.Join(t.TempDir(), "x", "open-error.db")}
	if _, err := NewSQLiteDB(context.Background(), cfg); err == nil {
		t.Fatalf("expected open error")
	}
}

func TestNewSQLiteDBPingError(t *testing.T) {
	origOpen := openSQLiteDB
	origPing := pingSQLiteDB
	t.Cleanup(func() {
		openSQLiteDB = origOpen
		pingSQLiteDB = origPing
	})

	db, err := sql.Open("sqlite", "file:"+filepath.Join(t.TempDir(), "ping.db"))
	if err != nil {
		t.Fatalf("open sqlite error: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close sqlite error: %v", err)
	}

	openSQLiteDB = func(string, string) (*sql.DB, error) { return db, nil }
	pingSQLiteDB = func(context.Context, *sql.DB) error { return errors.New("ping error") }

	cfg := &config.Config{SQLitePath: filepath.Join(t.TempDir(), "x", "ping-error.db")}
	if _, err := NewSQLiteDB(context.Background(), cfg); err == nil {
		t.Fatalf("expected ping error")
	}
}

func TestNewSQLiteDBEnsureSchemaError(t *testing.T) {
	origSchema := ensureSQLiteSchema
	origPing := pingSQLiteDB
	t.Cleanup(func() {
		ensureSQLiteSchema = origSchema
		pingSQLiteDB = origPing
	})

	ensureSQLiteSchema = func(context.Context, *sql.DB) error { return errors.New("schema error") }
	pingSQLiteDB = func(context.Context, *sql.DB) error { return nil }

	cfg := &config.Config{SQLitePath: filepath.Join(t.TempDir(), "x", "schema-error.db")}
	if _, err := NewSQLiteDB(context.Background(), cfg); err == nil {
		t.Fatalf("expected ensure schema error")
	}
}
