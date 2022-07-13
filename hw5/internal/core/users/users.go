package users

import (
	"context"
	"fmt"
)

// Repository holds datastore
type Repository interface {
	Create(ctx context.Context, u *User) error
	Read(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, u *User) error
}

// Users struct holds the Repository instance
type Users struct {
	repo Repository
}

// NewService creates a new Users service
func NewService(repo Repository) *Users {
	return &Users{repo: repo}
}

// NewModel create new instance User model
func (u Users) NewModel() *User {
	return &User{}
}

// Create save entity
func (u Users) Create(ctx context.Context, user *User) error {
	if err := u.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user in repo: %w", err)
	}

	return nil
}

// Read update entity
func (u Users) Read(ctx context.Context, user *User) error {
	if err := u.repo.Read(ctx, user); err != nil {
		return fmt.Errorf("failed to read user in repo: %w", err)
	}

	return nil
}

// Update update entity
func (u Users) Update(ctx context.Context, user *User) error {
	if err := u.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user in repo: %w", err)
	}

	return nil
}

// Delete update entity
func (u Users) Delete(ctx context.Context, user *User) error {
	if err := u.repo.Delete(ctx, user); err != nil {
		return fmt.Errorf("failed to delete user in repo: %w", err)
	}

	return nil
}
