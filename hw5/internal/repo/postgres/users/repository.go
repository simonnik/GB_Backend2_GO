package users

import (
	"context"
	"database/sql"

	"github.com/simonnik/GB_Backend2_GO/hw5/internal/core/users"
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/repo/datastore"
)

// UserRepository holds datastore
type UserRepository struct {
	Datastore *datastore.Datastore
}

// NewRepository create new instance repository
func NewRepository(d *datastore.Datastore) users.Repository {
	return &UserRepository{d}
}

// connection create connection
func (r *UserRepository) connection(u *users.User) (*sql.DB, error) {
	s, err := r.Datastore.M.ShardById(u.UserId)
	if err != nil {
		return nil, err
	}
	return r.Datastore.P.Connection(s.Address)
}

// Create user entity
func (r *UserRepository) Create(ctx context.Context, u *users.User) error {
	c, err := r.connection(u)
	if err != nil {
		return err
	}

	_, err = c.ExecContext(
		ctx,
		`INSERT INTO "users" VALUES ($1, $2, $3, $4)`,
		u.UserId,
		u.Name,
		u.Age,
		u.Spouse,
	)
	return err
}

// Read retrieve user
func (r *UserRepository) Read(ctx context.Context, u *users.User) error {
	c, err := r.connection(u)
	if err != nil {
		return err
	}

	q := c.QueryRowContext(ctx, `SELECT "name", "age", "spouse" FROM "users" WHERE "user_id" =
$1`, u.UserId)
	return q.Scan(
		&u.Name,
		&u.Age,
		&u.Spouse,
	)
}

// Update update user info
func (r *UserRepository) Update(ctx context.Context, u *users.User) error {
	c, err := r.connection(u)
	if err != nil {
		return err
	}

	_, err = c.ExecContext(ctx, `UPDATE "users" SET "name" = $2, "age" = $3, "spouse" = $4
WHERE "user_id" = $1`, u.UserId,
		u.Name, u.Age, u.Spouse)
	return err
}

// Delete user
func (r *UserRepository) Delete(ctx context.Context, u *users.User) error {
	c, err := r.connection(u)
	if err != nil {
		return err
	}

	_, err = c.ExecContext(ctx, `DELETE FROM "users" WHERE "user_id" = $1`, u.UserId)
	return err
}
