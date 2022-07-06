package datastore

import (
	"database/sql"
	"fmt"

	"github.com/simonnik/GB_Backend2_GO/hw3/internal/config"
)

type Datastore struct {
	DB *sql.DB
}

func NewDatastore(cfg *config.Config) (*Datastore, error) {
	connStr := "user=%s password=%s dbname=%s host=%s port=%d sslmode=%s "

	db, err := sql.Open("postgres", fmt.Sprintf(
		connStr,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.SSLMode,
	))

	if err != nil {
		return nil, err
	}

	return &Datastore{db}, nil
}

func (p Datastore) Close() error {
	return p.DB.Close()
}
