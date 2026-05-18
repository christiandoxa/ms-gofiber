package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	_ "modernc.org/sqlite"

	"ms-gofiber/internal/app/domain"
)

func setupDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+filepath.Join(t.TempDir(), "repo.db"))
	if err != nil {
		t.Fatalf("open sqlite error: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil && !errors.Is(err, sql.ErrConnDone) {
			t.Logf("close sqlite fixture: %v", err)
		}
	})
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS todos (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    completed INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos (created_at DESC);
`)
	if err != nil {
		t.Fatalf("create schema error: %v", err)
	}
	return db
}

func TestTodoRepositoryCRUD(t *testing.T) {
	db := setupDB(t)
	repo := NewTodo(db)
	ctx := context.Background()
	now := time.Now().UTC()

	createTodo(t, repo, ctx, "1", "t1", true, now)
	createTodo(t, repo, ctx, "2", "t2", false, now.Add(time.Second))

	got, err := repo.GetByID(ctx, "1")
	if err != nil {
		t.Fatalf("get error: %v", err)
	}
	if got.ID != "1" || !got.Completed {
		t.Fatalf("unexpected todo: %+v", got)
	}

	if _, err := repo.GetByID(ctx, "404"); !errors.Is(err, domain.ErrTodoNotFound) {
		t.Fatalf("expected ErrTodoNotFound, got %v", err)
	}

	list, err := repo.List(ctx, -1, -1)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(list) < 2 {
		t.Fatalf("expected at least 2 rows, got %d", len(list))
	}

	if err := repo.Update(ctx, &domain.Todo{ID: "1", Title: "updated", Completed: false, UpdatedAt: now.Add(time.Hour)}); err != nil {
		t.Fatalf("update error: %v", err)
	}
	if err := repo.Update(ctx, &domain.Todo{ID: "none", Title: "x", UpdatedAt: now}); !errors.Is(err, domain.ErrTodoNotFound) {
		t.Fatalf("expected update not found, got %v", err)
	}

	if err := repo.Delete(ctx, "2"); err != nil {
		t.Fatalf("delete error: %v", err)
	}
	if err := repo.Delete(ctx, "none"); !errors.Is(err, domain.ErrTodoNotFound) {
		t.Fatalf("expected delete not found, got %v", err)
	}
}

func TestTodoRepositoryTimestampBranches(t *testing.T) {
	db := setupDB(t)
	repo := NewTodo(db)
	ctx := context.Background()

	_, err := db.Exec(`INSERT INTO todos(id,title,completed,created_at,updated_at) VALUES(?,?,?,?,?)`, "3", "t3", 1, time.Now().UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("insert fallback row error: %v", err)
	}
	if _, err = repo.GetByID(ctx, "3"); err != nil {
		t.Fatalf("fallback parse get error: %v", err)
	}

	_, err = db.Exec(`INSERT INTO todos(id,title,completed,created_at,updated_at) VALUES(?,?,?,?,?)`, "bad-ts", "x", 0, "bad", "bad")
	if err != nil {
		t.Fatalf("insert bad row error: %v", err)
	}
	if _, err = repo.GetByID(ctx, "bad-ts"); err == nil {
		t.Fatalf("expected parse error for bad timestamp")
	}
}

func TestTodoRepositoryClosedDBBranches(t *testing.T) {
	db := setupDB(t)
	repo := NewTodo(db)
	ctx := context.Background()
	now := time.Now().UTC()

	if err := db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}
	if _, err := repo.List(ctx, 1, 0); err == nil {
		t.Fatalf("expected list error on closed db")
	}
	if err := repo.Update(ctx, &domain.Todo{ID: "1", Title: "x", UpdatedAt: now}); err == nil {
		t.Fatalf("expected update error on closed db")
	}
	if err := repo.Delete(ctx, "1"); err == nil {
		t.Fatalf("expected delete error on closed db")
	}
}

func TestTodoRepositoryListCloseError(t *testing.T) {
	db := setupDB(t)
	repo := NewTodo(db)
	ctx := context.Background()
	now := time.Now().UTC()
	createTodo(t, repo, ctx, "close-error", "x", false, now)

	expected := errors.New("close rows error")
	var rows *sql.Rows
	patches := gomonkey.ApplyMethod(rows, "Close", func(*sql.Rows) error {
		return expected
	})
	defer patches.Reset()

	if _, err := repo.List(ctx, 1, 0); !errors.Is(err, expected) {
		t.Fatalf("expected close rows error, got %v", err)
	}
}

func createTodo(
	t *testing.T,
	repo *Todo,
	ctx context.Context,
	id domain.TodoID,
	title string,
	completed bool,
	now time.Time,
) {
	t.Helper()

	_, err := repo.Create(ctx, &domain.Todo{
		ID:        id,
		Title:     title,
		Completed: completed,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		t.Fatalf("create error: %v", err)
	}
}

func TestParseAndBoolHelpers(t *testing.T) {
	if got := boolToInt(true); got != 1 {
		t.Fatalf("expected 1 got %d", got)
	}
	if got := boolToInt(false); got != 0 {
		t.Fatalf("expected 0 got %d", got)
	}
	if _, err := parseRFC3339("2026-02-23T10:30:00.123456789Z"); err != nil {
		t.Fatalf("nano parse error: %v", err)
	}
	if _, err := parseRFC3339("2026-02-23T10:30:00Z"); err != nil {
		t.Fatalf("fallback parse error: %v", err)
	}
	if _, err := parseRFC3339("bad"); err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestTodoRepositoryAdditionalBranches(t *testing.T) {
	db := setupDB(t)
	repo := NewTodo(db)
	ctx := context.Background()
	now := time.Now().UTC()

	// list scan error branch
	if _, err := db.Exec(
		`INSERT INTO todos(id,title,completed,created_at,updated_at) VALUES(?,?,?,?,?)`,
		"bad-completed",
		"x",
		"not-int",
		now.Format(time.RFC3339Nano),
		now.Format(time.RFC3339Nano),
	); err != nil {
		t.Fatalf("insert bad completed row error: %v", err)
	}
	if _, err := repo.List(ctx, 1, 0); err == nil {
		t.Fatalf("expected list scan error")
	}

	// scanTodo updatedAt parse error branch
	if _, err := db.Exec(
		`INSERT INTO todos(id,title,completed,created_at,updated_at) VALUES(?,?,?,?,?)`,
		"bad-updated",
		"x",
		1,
		now.Format(time.RFC3339Nano),
		"bad-updated-time",
	); err != nil {
		t.Fatalf("insert bad updated row error: %v", err)
	}
	if _, err := repo.GetByID(ctx, "bad-updated"); err == nil {
		t.Fatalf("expected updated_at parse error")
	}

	if _, err := repo.Create(ctx, &domain.Todo{
		ID:        "rows-affected",
		Title:     "x",
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}); err != nil {
		t.Fatalf("create rows-affected row error: %v", err)
	}

	sampleResult, err := db.Exec(`UPDATE todos SET title = title WHERE id = ?`, "rows-affected")
	if err != nil {
		t.Fatalf("create rows affected sample result: %v", err)
	}
	patches := gomonkey.ApplyMethod(reflect.TypeOf(sampleResult), "RowsAffected", func(sql.Result) (int64, error) {
		return 0, errors.New("rows affected error")
	})
	defer patches.Reset()

	if err := repo.Update(ctx, &domain.Todo{
		ID:        "rows-affected",
		Title:     "u",
		Completed: true,
		UpdatedAt: now.Add(time.Minute),
	}); err == nil {
		t.Fatalf("expected rows affected error on update")
	}

	if err := repo.Delete(ctx, "rows-affected"); err == nil {
		t.Fatalf("expected rows affected error on delete")
	}
}
