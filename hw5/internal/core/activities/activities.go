package activities

import (
	"context"
	"fmt"
)

// Repository holds datastore
type Repository interface {
	Create(ctx context.Context, a *Activity) error
}

// Activities struct holds the Repository instance
type Activities struct {
	repo Repository
}

// NewService creates a new Activities service
func NewService(repo Repository) *Activities {
	return &Activities{repo: repo}
}

// Create save entity
func (a Activities) Create(ctx context.Context, act *Activity) error {
	if err := a.repo.Create(ctx, act); err != nil {
		return fmt.Errorf("failed to create activity in repo: %w", err)
	}

	return nil
}
