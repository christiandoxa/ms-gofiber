package database

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
)

func TestTodoDatabaseCreateGetList(t *testing.T) {
	ctx := context.Background()
	db, now := newTodoDatabase(t)

	if _, err := db.CreateTodo(ctx, TodoRecord{ID: "2", Title: "second", CreatedAt: now.Add(time.Second), UpdatedAt: now.Add(time.Second)}); err != nil {
		t.Fatalf("create second: %v", err)
	}
	if _, err := db.CreateTodo(ctx, TodoRecord{ID: "1", Title: "first", CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatalf("create first: %v", err)
	}

	record, ok, err := db.GetTodo(ctx, "1")
	if err != nil || !ok || record.Title != "first" {
		t.Fatalf("unexpected get result: %+v %v %v", record, ok, err)
	}
	if _, ok, err := db.GetTodo(ctx, "missing"); err != nil || ok {
		t.Fatalf("unexpected missing get result: %v %v", ok, err)
	}

	records, err := db.ListTodos(ctx)
	if err != nil || len(records) != 2 || records[0].ID != "1" {
		t.Fatalf("unexpected list result: %+v %v", records, err)
	}
}

func TestTodoDatabaseUpdateDelete(t *testing.T) {
	ctx := context.Background()
	db, now := newTodoDatabase(t)

	if _, err := db.CreateTodo(ctx, TodoRecord{ID: "1", Title: "first", CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatalf("create first: %v", err)
	}

	updated, ok, err := db.UpdateTodo(ctx, TodoRecord{ID: "1", Title: "updated", CreatedAt: now, UpdatedAt: now})
	if err != nil || !ok || updated.Title != "updated" {
		t.Fatalf("unexpected update result: %+v %v %v", updated, ok, err)
	}
	if _, ok, err := db.UpdateTodo(ctx, TodoRecord{ID: "missing"}); err != nil || ok {
		t.Fatalf("unexpected missing update result: %v %v", ok, err)
	}
	if ok, err := db.DeleteTodo(ctx, "1"); err != nil || !ok {
		t.Fatalf("delete failed: %v %v", ok, err)
	}
	if ok, err := db.DeleteTodo(ctx, "missing"); err != nil || ok {
		t.Fatalf("unexpected missing delete result: %v %v", ok, err)
	}
}

func newTodoDatabase(t *testing.T) (*DB, time.Time) {
	t.Helper()
	db := Connect(":memory:")
	t.Cleanup(func() {
		closeDB(t, db)
	})
	return db, time.Now().UTC()
}

func TestOpenFileDatabase(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "nested", "todos.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	closeDB(t, db)
}

func TestOpenError(t *testing.T) {
	if _, err := Open("/proc/1/forbidden/todos.db"); err == nil {
		t.Fatalf("expected open error")
	}
}

func TestOpenSQLError(t *testing.T) {
	patches := gomonkey.ApplyGlobalVar(&openSQL, func(string, string) (*sql.DB, error) {
		return nil, errors.New("open")
	})
	defer patches.Reset()
	if _, err := Open(":memory:"); err == nil {
		t.Fatalf("expected sql open error")
	}
}

func TestOpenSchemaError(t *testing.T) {
	if _, err := Open("/"); err == nil {
		t.Fatalf("expected schema error")
	}
}

func TestConnectPanic(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("expected panic")
		}
	}()
	Connect("/proc/1/forbidden/todos.db")
}

func TestClosedDatabaseErrors(t *testing.T) {
	ctx := context.Background()
	db := Connect(":memory:")
	closeDB(t, db)
	record := TodoRecord{ID: "1", Title: "title", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	if _, err := db.CreateTodo(ctx, record); err == nil {
		t.Fatalf("expected create error")
	}
	if _, _, err := db.GetTodo(ctx, "1"); err == nil {
		t.Fatalf("expected get error")
	}
	if _, err := db.ListTodos(ctx); err == nil {
		t.Fatalf("expected list error")
	}
	if _, _, err := db.UpdateTodo(ctx, record); err == nil {
		t.Fatalf("expected update error")
	}
	if _, err := db.DeleteTodo(ctx, "1"); err == nil {
		t.Fatalf("expected delete error")
	}
}

func TestScanTodoParseError(t *testing.T) {
	db := Connect(":memory:")
	defer closeDB(t, db)
	_, err := db.sqlDB.ExecContext(
		context.Background(),
		`INSERT INTO todos (id, title, completed, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		"1",
		"title",
		1,
		"bad",
		time.Now().UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		t.Fatalf("insert bad record: %v", err)
	}
	if _, _, err := db.GetTodo(context.Background(), "1"); err == nil {
		t.Fatalf("expected created_at parse error")
	}
}

func TestListTodosScanError(t *testing.T) {
	db := Connect(":memory:")
	defer closeDB(t, db)
	_, err := db.sqlDB.ExecContext(
		context.Background(),
		`INSERT INTO todos (id, title, completed, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		"1",
		"title",
		1,
		"bad",
		time.Now().UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		t.Fatalf("insert bad record: %v", err)
	}
	if _, err := db.ListTodos(context.Background()); err == nil {
		t.Fatalf("expected list parse error")
	}
}

func TestRowsAffectedErrors(t *testing.T) {
	db := Connect(":memory:")
	defer closeDB(t, db)
	record := TodoRecord{ID: "1", Title: "title", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	if _, err := db.CreateTodo(context.Background(), record); err != nil {
		t.Fatalf("create todo: %v", err)
	}

	patches := gomonkey.ApplyGlobalVar(&getRowsAffected, func(sql.Result) (int64, error) {
		return 0, errors.New("affected")
	})
	defer patches.Reset()

	if _, _, err := db.UpdateTodo(context.Background(), record); err == nil {
		t.Fatalf("expected update affected error")
	}
	if _, err := db.DeleteTodo(context.Background(), "1"); err == nil {
		t.Fatalf("expected delete affected error")
	}
}

func TestScanTodoUpdatedAtParseError(t *testing.T) {
	db := Connect(":memory:")
	defer closeDB(t, db)
	_, err := db.sqlDB.ExecContext(
		context.Background(),
		`INSERT INTO todos (id, title, completed, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		"1",
		"title",
		1,
		time.Now().UTC().Format(time.RFC3339Nano),
		"bad",
	)
	if err != nil {
		t.Fatalf("insert bad record: %v", err)
	}
	if _, _, err := db.GetTodo(context.Background(), "1"); err == nil {
		t.Fatalf("expected updated_at parse error")
	}
}

func TestBoolToInt(t *testing.T) {
	if boolToInt(true) != 1 {
		t.Fatalf("expected true as 1")
	}
	if boolToInt(false) != 0 {
		t.Fatalf("expected false as 0")
	}
}

func closeDB(t *testing.T, db *DB) {
	t.Helper()
	if err := db.Close(); err != nil {
		t.Fatalf("close database: %v", err)
	}
}
