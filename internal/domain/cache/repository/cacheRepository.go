package repository

import (
	"context"

	"ms-gofiber/pkg/infrastructure/cache"
)

type ICacheRepository interface {
	Flush(ctx context.Context) error
}

type CacheRepository struct {
	cache *cache.Cache
}

func New(cacheClient *cache.Cache) ICacheRepository {
	return &CacheRepository{cache: cacheClient}
}

func (r *CacheRepository) Flush(ctx context.Context) error {
	return r.cache.Flush(ctx)
}
