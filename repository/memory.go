// FOR TESTING AND LOCAL USE ONLY
// DOESNT HAVE EXPIRY CLEANUP

package repository

import (
	"context"
	"sync"
)

type MemoryRepo struct {
	mu    sync.RWMutex
	store map[string]int
}

func NewMemoryRepo() RepoInterface {
	return &MemoryRepo{store: make(map[string]int)}
}

func (m *MemoryRepo) Increment(ctx context.Context, key string, expiry int) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[key] += 1
	return m.store[key], nil
}

func (m *MemoryRepo) GetByKey(ctx context.Context, key string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := m.store[key]
	return result, nil
}
