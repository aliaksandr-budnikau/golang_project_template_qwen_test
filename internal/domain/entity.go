// Package domain contains the core business entities and the ports (interfaces)
// the application layer depends on. It has no dependencies on any other
// internal package.
package domain

import (
	"context"
	"errors"
)

// ErrNotFound is returned by repositories when an item does not exist.
var ErrNotFound = errors.New("domain: item not found")

// Item is a sample domain entity.
type Item struct {
	ID   string
	Name string
}

// Repository is the persistence port the application layer depends on.
// Concrete implementations live in the infrastructure layer.
type Repository interface {
	// Save persists the item, overwriting any existing one with the same id.
	Save(ctx context.Context, item Item) error
	// Get returns the item with the given id, or ErrNotFound.
	Get(ctx context.Context, id string) (Item, error)
}
