package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/simonnik/GB_Backend2_GO/internal/entity"
)

// Repo is interface to interact with particular storage repository
type Repo interface {
	CreateUser(ctx context.Context, u entity.User) error
	CreateGroup(ctx context.Context, u entity.Group) error
	AddToGroup(ctx context.Context, uid, gid uuid.UUID) error
	RemoveFromGroup(ctx context.Context, uid, gid uuid.UUID) error
	SearchUser(ctx context.Context, name string, gids []uuid.UUID) ([]entity.User, error)
	SearchGroup(ctx context.Context, name string, uids []uuid.UUID) ([]entity.Group, error)
}

// DB implements main logic using Repo methods
type DB struct {
	repo Repo
}

// NewDB creates Storage instance
func NewDB(repo Repo) *DB {
	return &DB{
		repo: repo,
	}
}

// CreateUser adds ID to the passed user data and calls repo's CreateUser method
func (db *DB) CreateUser(ctx context.Context, u entity.User) (*entity.User, error) {
	u.ID = uuid.New()
	err := db.repo.CreateUser(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("error creating new user: %w", err)
	}
	return &u, nil
}

// CreateGroup adds ID to the passed group data and calls repo's CreateGroup method
func (db *DB) CreateGroup(ctx context.Context, g entity.Group) (*entity.Group, error) {
	g.ID = uuid.New()
	err := db.repo.CreateGroup(ctx, g)
	if err != nil {
		return nil, fmt.Errorf("error creating new user: %w", err)
	}
	return &g, nil
}

// AddToGroup calls same function from repo
func (db *DB) AddToGroup(ctx context.Context, uid, gid uuid.UUID) error {
	return db.repo.AddToGroup(ctx, uid, gid)
}

// RemoveFromGroup calls same function from repo
func (db *DB) RemoveFromGroup(ctx context.Context, uid, gid uuid.UUID) error {
	return db.repo.RemoveFromGroup(ctx, uid, gid)
}

// SearchUser ...
func (db *DB) SearchUser(ctx context.Context, name string, gids []uuid.UUID) ([]entity.User, error) {

	if name == "" && len(gids) == 0 {
		return nil, errors.New("name and gids are empty, at least one criteria should be set")
	}

	users, err := db.repo.SearchUser(ctx, name, gids)
	if err != nil {
		return nil, fmt.Errorf("error searching users: %w", err)
	}
	return users, nil
}

// SearchGroup ...
func (db *DB) SearchGroup(ctx context.Context, name string, uids []uuid.UUID) ([]entity.Group, error) {

	if name == "" && len(uids) == 0 {
		return nil, errors.New("name and uids are empty, at least one criteria should be set")
	}

	groups, err := db.repo.SearchGroup(ctx, name, uids)
	if err != nil {
		return nil, fmt.Errorf("error searching groups: %w", err)
	}
	return groups, nil
}
