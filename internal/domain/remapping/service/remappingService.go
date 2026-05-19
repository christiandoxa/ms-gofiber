package service

import (
	"context"

	"ms-gofiber/internal/domain/remapping/repository"
	"ms-gofiber/pkg/responsecode/model"
)

type IRemappingService interface {
	Remapping(ctx context.Context, endpoint string, originResponseCode string, originResponseMessage string) (*rcmodel.ResponseCode, error)
}

type RemappingService struct {
	remappingRepository repository.IRemappingRepository
}

func New(remappingRepository repository.IRemappingRepository) IRemappingService {
	return &RemappingService{remappingRepository: remappingRepository}
}

func (s *RemappingService) Remapping(ctx context.Context, endpoint string, originResponseCode string, originResponseMessage string) (*rcmodel.ResponseCode, error) {
	responseCode, err := s.remappingRepository.Remapping(ctx, endpoint, originResponseCode, originResponseMessage)
	if err != nil {
		return nil, err
	}
	return rcmodel.NewResponseCode(responseCode.StatusCode, responseCode.ResponseMessage, responseCode.ResponseCode), nil
}
