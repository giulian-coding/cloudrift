package containers

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	mu         sync.RWMutex
	containers map[string]Container
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{containers: make(map[string]Container)}
}

func (r *MemoryRepository) Create(ctx context.Context, c Container) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.containers[c.ID] = c
	return nil
}

func (r *MemoryRepository) ByID(ctx context.Context, id string) (Container, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if order, exists := r.containers[id]; exists {
		return order, nil
	}

	return Container{}, ErrNotFound
}
