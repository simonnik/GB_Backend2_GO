package activities

import (
	"context"
	"database/sql"

	"github.com/simonnik/GB_Backend2_GO/hw5/internal/core/activities"
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/repo/datastore"
)

// ActivityRepository holds datastore
type ActivityRepository struct {
	Datastore *datastore.Datastore
}

// NewRepository create new instance repository
func NewRepository(d *datastore.Datastore) activities.Repository {
	return &ActivityRepository{d}
}

// connection create connection
func (r *ActivityRepository) connection(userId int) (*sql.DB, error) {
	s, err := r.Datastore.M.ShardById(userId)
	if err != nil {
		return nil, err
	}
	return r.Datastore.P.Connection(s.Address)
}

// Create user entity
func (r *ActivityRepository) Create(ctx context.Context, a *activities.Activity) error {
	c, err := r.connection(a.UserId)
	if err != nil {
		return err
	}

	_, err = c.ExecContext(
		ctx,
		`INSERT INTO "activities" (user_id, name, date) VALUES ($1, $2, $3)`,
		a.UserId,
		a.Name,
		a.Date,
	)
	return err
}
