// Package application contains the use cases (application services) that
// orchestrate the domain to fulfil individual requests. It depends on the
// domain ports, never on their concrete implementations.
package application

import (
	"context"

	"github.com/yourname/empty_project/internal/domain"
)

// UseCase is a sample application service. It depends on the domain.Repository
// port, not on any concrete implementation.
type UseCase struct {
	repo domain.Repository
}

// NewUseCase constructs a UseCase.
func NewUseCase(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Create persists a new item.
func (uc *UseCase) Create(ctx context.Context, item domain.Item) error {
	return uc.repo.Save(ctx, item)
}

// Retrieve loads an item by id.
func (uc *UseCase) Retrieve(ctx context.Context, id string) (domain.Item, error) {
	return uc.repo.Get(ctx, id)
}
