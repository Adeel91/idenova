package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/Adeel91/idenova/services/verification/internal/domain"
)

type Repository struct {
	mu   sync.Mutex
	data map[string]*domain.Verification
}

func NewRepository() *Repository {
	return &Repository{
		data: make(map[string]*domain.Verification),
	}
}

func (r *Repository) Save(ctx context.Context, v *domain.Verification) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	copy := *v
	r.data[v.ID] = &copy
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*domain.Verification, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	v, exists := r.data[id]
	if !exists {
		return nil, errors.New("verification not found")
	}

	copy := *v
	return &copy, nil
}
