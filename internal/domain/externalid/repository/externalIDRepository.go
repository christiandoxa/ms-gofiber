package repository

import (
	"context"
	"errors"

	"ms-gofiber/pkg/constant/generalkey"
	"ms-gofiber/pkg/infrastructure/cache"
	"ms-gofiber/pkg/responsecode/error"
	"ms-gofiber/pkg/tool"
)

type IExternalIDRepository interface {
	StoreExternalID(ctx context.Context, externalID string) error
}

type ExternalIDRepository struct {
	cache *cache.Cache
}

func New(cacheClient *cache.Cache) IExternalIDRepository {
	return &ExternalIDRepository{cache: cacheClient}
}

func (r *ExternalIDRepository) StoreExternalID(ctx context.Context, externalID string) error {
	err := r.cache.Store(ctx, cache.Data{
		Key:      generalkey.ExternalIDPrefix + externalID,
		Value:    true,
		Duration: tool.GetExpiration(),
	})
	if errors.Is(err, cache.ErrCacheKeyExists) {
		return rcerror.ErrDuplicateExternalID
	}
	return err
}
