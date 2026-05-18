package database

import (
	"testing"
	"time"
)

func TestTodoDatabase(t *testing.T) {
	db := Connect()
	now := time.Now().UTC()

	db.CreateTodo(TodoRecord{ID: "2", Title: "second", CreatedAt: now.Add(time.Second), UpdatedAt: now.Add(time.Second)})
	db.CreateTodo(TodoRecord{ID: "1", Title: "first", CreatedAt: now, UpdatedAt: now})

	record, ok := db.GetTodo("1")
	if !ok || record.Title != "first" {
		t.Fatalf("unexpected get result: %+v %v", record, ok)
	}

	if records := db.ListTodos(); len(records) != 2 || records[0].ID != "1" {
		t.Fatalf("unexpected list result: %+v", records)
	}

	updated, ok := db.UpdateTodo(TodoRecord{ID: "1", Title: "updated", CreatedAt: now, UpdatedAt: now})
	if !ok || updated.Title != "updated" {
		t.Fatalf("unexpected update result: %+v %v", updated, ok)
	}

	if _, ok := db.UpdateTodo(TodoRecord{ID: "missing"}); ok {
		t.Fatalf("expected missing update")
	}
	if !db.DeleteTodo("1") {
		t.Fatalf("expected delete success")
	}
	if db.DeleteTodo("missing") {
		t.Fatalf("expected missing delete")
	}
}
