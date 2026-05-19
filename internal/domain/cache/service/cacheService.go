package service

import (
	"context"

	"ms-gofiber/internal/domain/cache/repository"
)

type ICacheService interface {
	Flush(ctx context.Context) error
}

type CacheService struct {
	cacheRepository repository.ICacheRepository
}

func New(cacheRepository repository.ICacheRepository) ICacheService {
	return &CacheService{cacheRepository: cacheRepository}
}

func (s *CacheService) Flush(ctx context.Context) error {
	return s.cacheRepository.Flush(ctx)
}
