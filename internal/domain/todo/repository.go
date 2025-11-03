package todo

import "context"

type Repository interface {
	Create(ctx context.Context, t *Todo) (ID, error)
	GetByID(ctx context.Context, id ID) (*Todo, error)
	List(ctx context.Context, limit, offset int) ([]*Todo, error)
	Update(ctx context.Context, t *Todo) error
	Delete(ctx context.Context, id ID) error
}
