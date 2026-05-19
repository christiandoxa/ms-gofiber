package service

import (
	"context"

	"ms-gofiber/internal/domain/externalid/repository"
)

type IExternalIDService interface {
	StoreExternalID(ctx context.Context, externalID string) error
}

type ExternalIDService struct {
	externalIDRepository repository.IExternalIDRepository
}

func New(externalIDRepository repository.IExternalIDRepository) IExternalIDService {
	return &ExternalIDService{externalIDRepository: externalIDRepository}
}

func (s *ExternalIDService) StoreExternalID(ctx context.Context, externalID string) error {
	return s.externalIDRepository.StoreExternalID(ctx, externalID)
}
