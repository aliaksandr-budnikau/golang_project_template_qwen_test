// Package persistence provides concrete implementations of the domain ports for
// persistence. It depends on the domain, never the other way round.
package persistence

import (
	"context"
	"sync"

	"github.com/yourname/empty_project/internal/domain"
)

// MemoryRepo is an in-memory implementation of domain.Repository. It is safe for
// concurrent use and is intended for local development and tests.
type MemoryRepo struct {
	mu    sync.RWMutex
	items map[string]domain.Item
}

// NewMemoryRepo constructs an empty MemoryRepo.
func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{items: make(map[string]domain.Item)}
}

// Save stores an item, overwriting any existing one with the same id.
func (r *MemoryRepo) Save(_ context.Context, item domain.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[item.ID] = item
	return nil
}

// Get returns the item with the given id, or domain.ErrNotFound.
func (r *MemoryRepo) Get(_ context.Context, id string) (domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[id]
	if !ok {
		return domain.Item{}, domain.ErrNotFound
	}
	return item, nil
}

// Compile-time assertion that *MemoryRepo satisfies domain.Repository.
var _ domain.Repository = (*MemoryRepo)(nil)
