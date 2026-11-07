package application

import (
	"context"
	"testing"

	"github.com/yourname/empty_project/internal/domain"
)

// fakeRepo is an in-memory domain.Repository for tests.
type fakeRepo struct {
	items map[string]domain.Item
}

func (f *fakeRepo) Save(_ context.Context, item domain.Item) error {
	f.items[item.ID] = item
	return nil
}

func (f *fakeRepo) Get(_ context.Context, id string) (domain.Item, error) {
	item, ok := f.items[id]
	if !ok {
		return domain.Item{}, domain.ErrNotFound
	}
	return item, nil
}

func TestUseCaseCreateThenRetrieve(t *testing.T) {
	uc := NewUseCase(&fakeRepo{items: make(map[string]domain.Item)})
	item := domain.Item{ID: "1", Name: "alpha"}

	if err := uc.Create(context.Background(), item); err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	got, err := uc.Retrieve(context.Background(), item.ID)
	if err != nil {
		t.Fatalf("Retrieve() returned error: %v", err)
	}
	if got != item {
		t.Errorf("Retrieve() = %+v, want %+v", got, item)
	}
}

func TestUseCaseRetrieveMissing(t *testing.T) {
	uc := NewUseCase(&fakeRepo{items: make(map[string]domain.Item)})

	if _, err := uc.Retrieve(context.Background(), "nope"); err != domain.ErrNotFound {
		t.Errorf("Retrieve() error = %v, want %v", err, domain.ErrNotFound)
	}
}
