package presenter

import (
	"testing"
	"time"

	"ms-gofiber/internal/app/domain"
)

func TestTodoData(t *testing.T) {
	if got := TodoData(nil); got != (Todo{}) {
		t.Fatalf("expected zero value for nil input: %+v", got)
	}

	now := time.Now().UTC()
	in := &domain.Todo{ID: "1", Title: "x", Completed: true, CreatedAt: now, UpdatedAt: now}
	got := TodoData(in)
	if got.ID != "1" || got.Title != "x" || !got.Completed {
		t.Fatalf("unexpected presenter mapping: %+v", got)
	}

	list := TodoListData([]*domain.Todo{in})
	if len(list) != 1 || list[0].ID != "1" {
		t.Fatalf("unexpected list presenter mapping: %+v", list)
	}
}
