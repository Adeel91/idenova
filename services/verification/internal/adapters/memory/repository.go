package memory

import (
	"context"
	"errors"
	"sort"
	"sync"

	"idenova/verification/internal/domain"
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

func (r *Repository) Update(ctx context.Context, v *domain.Verification) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[v.ID]; !exists {
		return errors.New("verification not found")
	}

	copy := *v
	r.data[v.ID] = &copy

	return nil
}

func (r *Repository) List(
	ctx context.Context,
	status *domain.VerificationStatus,
	limit, offset int32,
) ([]domain.Verification, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if limit <= 0 || offset < 0 {
		return nil, errors.New("invalid pagination parameters")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	results := make([]domain.Verification, 0)

	for _, v := range r.data {
		if status == nil || v.Status == *status {
			results = append(results, *v)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].ID > results[j].ID
	})

	if int64(offset) >= int64(len(results)) {
		return []domain.Verification{}, nil
	}

	end := int64(offset) + int64(limit)
	if end > int64(len(results)) {
		end = int64(len(results))
	}

	return results[int(offset):int(end)], nil
}
