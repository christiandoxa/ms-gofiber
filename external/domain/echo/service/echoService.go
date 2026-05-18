package service

import (
	"context"

	"ms-gofiber/external/domain/echo/repository"
)

type IEchoService interface {
	Echo(ctx context.Context, target string) (string, error)
}

type EchoService struct {
	echoRepository repository.IEchoRepository
}

func New(echoRepository repository.IEchoRepository) IEchoService {
	return &EchoService{echoRepository: echoRepository}
}

func (s *EchoService) Echo(ctx context.Context, target string) (string, error) {
	response, err := s.echoRepository.Fetch(ctx, target)
	if err != nil {
		return "", err
	}
	return string(response.Body), nil
}
