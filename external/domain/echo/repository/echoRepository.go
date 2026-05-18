package repository

import (
	"context"
	"net/http"

	"ms-gofiber/pkg/client"
)

type IEchoRepository interface {
	Fetch(ctx context.Context, target string) (*client.Response, error)
}

type EchoRepository struct{}

func New() IEchoRepository {
	return &EchoRepository{}
}

func (r *EchoRepository) Fetch(ctx context.Context, target string) (*client.Response, error) {
	return client.Do(ctx, client.Request{
		Method: http.MethodGet,
		URL:    target,
	})
}
