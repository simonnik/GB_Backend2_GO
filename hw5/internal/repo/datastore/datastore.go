package datastore

import (
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/pool"
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/sharding"
)

type Datastore struct {
	M *sharding.Manager
	P *pool.Pool
}

func NewDatastore(m *sharding.Manager, p *pool.Pool) (*Datastore, error) {
	return &Datastore{m, p}, nil
}
