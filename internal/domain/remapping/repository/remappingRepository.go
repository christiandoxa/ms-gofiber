package repository

import (
	"context"
	"errors"
	"sync"

	remappingmodel "ms-gofiber/internal/domain/remapping/model"
)

var ErrRemappingNotFound = errors.New("remapping not found")

type IRemappingRepository interface {
	Remapping(ctx context.Context, endpoint string, originResponseCode string, originResponseMessage string) (*remappingmodel.ResponseCode, error)
}

type RemappingRepository struct {
	items map[string]remappingmodel.ResponseCode
	mutex sync.RWMutex
}

func New(items ...remappingmodel.ResponseCode) IRemappingRepository {
	repository := &RemappingRepository{
		items: map[string]remappingmodel.ResponseCode{},
	}
	for _, item := range items {
		repository.items[cacheKey(item.Endpoint, item.OriginResponseCode, item.OriginResponseMessage)] = item
	}
	return repository
}

func (r *RemappingRepository) Remapping(_ context.Context, endpoint string, originResponseCode string, originResponseMessage string) (*remappingmodel.ResponseCode, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	responseCode, ok := r.items[cacheKey(endpoint, originResponseCode, originResponseMessage)]
	if !ok {
		return nil, ErrRemappingNotFound
	}
	return &responseCode, nil
}

func cacheKey(endpoint string, originResponseCode string, originResponseMessage string) string {
	return endpoint + "|" + originResponseCode + "|" + originResponseMessage
}
