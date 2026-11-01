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
	if _, exists := r.containers[c.ID]; !exists {
		r.containers[c.ID] = c
		return nil
	}
	return ErrAlreadyExists
}

func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.containers[id]; exists {
		delete(r.containers, id)
		return nil
	}

	return ErrNotFound
}

func (r *MemoryRepository) ByID(ctx context.Context, id string) (Container, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if container, exists := r.containers[id]; exists {
		return container, nil
	}

	return Container{}, ErrNotFound
}
