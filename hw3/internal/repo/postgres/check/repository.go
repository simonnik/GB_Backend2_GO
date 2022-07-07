package check

import (
	"fmt"

	"github.com/simonnik/GB_Backend2_GO/hw3/internal/core/check"
	"github.com/simonnik/GB_Backend2_GO/hw3/internal/repo/datastore"
)

// checkRepository holds datastore
type checkRepository struct {
	ds *datastore.Datastore
}

func NewRepository(d *datastore.Datastore) check.Repository {
	return &checkRepository{d}
}

func (c checkRepository) Check() error {
	row := c.ds.DB.QueryRow("SELECT 1 as ok")

	var ok string
	if err := row.Scan(&ok); err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	return nil
}
