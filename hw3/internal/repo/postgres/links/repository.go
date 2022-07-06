package links

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/neko-neko/echo-logrus/v2/log"

	"github.com/simonnik/GB_Backend2_GO/hw3/internal/core/links"
	"github.com/simonnik/GB_Backend2_GO/hw3/internal/repo/datastore"
)

// LinkRepository holds datastore
type LinkRepository struct {
	Datastore *datastore.Datastore
}

func NewRepository(d *datastore.Datastore) links.Repository {
	return &LinkRepository{d}
}

// Create link entity
func (r LinkRepository) Create(ctx context.Context, link links.Link) error {
	query := "INSERT INTO links (link, token) VALUES " +
		"($1, $2) RETURNING id"
	rows, err := r.Datastore.DB.QueryContext(ctx,
		query,
		link.Link,
		link.Token,
	)
	if err != nil {
		return fmt.Errorf("failed to insert link to db: %w", err)
	}

	rows.Next()
	err = rows.Scan(&link.ID)
	if err != nil {
		return fmt.Errorf("failed to get last inserted id from db: %w", err)
	}

	return nil
}

// FindByToken link entity
func (r LinkRepository) FindByToken(ctx context.Context, token string) (*links.Link, error) {
	l := links.Link{}
	query := "SELECT id, token, link FROM links WHERE token = $1"
	row := r.Datastore.DB.QueryRowContext(ctx, query, token)

	if err := row.Scan(
		&l.ID,
		&l.Token,
		&l.Link,
	); err != nil {
		return nil, fmt.Errorf("failed to select link from db: %w", err)
	}

	return &l, nil
}

// FindAllByToken link entity
func (r LinkRepository) FindAllByToken(ctx context.Context, token string) (links.StatList, error) {
	query := strings.Builder{}
	query.WriteString("SELECT l.id, l.link, ls.ip, ls.created_at FROM links_stat ls")
	query.WriteString(" JOIN links l ON l.id = ls.link_id")
	query.WriteString(" WHERE l.token = $1")

	stmt, err := r.Datastore.DB.PrepareContext(ctx, query.String())
	if err != nil {
		return nil, fmt.Errorf("faield to prepare statement: %w", err)
	}

	rows, err := stmt.QueryContext(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to select link from db: %w", err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Error(fmt.Errorf("error close rows: %w", err))
		}
	}()
	var stats links.StatList

	for rows.Next() {
		s := &links.Stat{}
		if err := rows.Scan(
			&s.ID,
			&s.IP,
			&s.Link,
			&s.Created,
		); err != nil {
			return nil, err
		}

		stats = append(stats, s)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to select link from db: %w", err)
	}
	return stats, nil
}

// SaveStat link entity
func (r LinkRepository) SaveStat(ctx context.Context, id int64, ip string) error {
	query := "INSERT INTO links_stat (link_id, ip) VALUES ($1, $2)"
	res, err := r.Datastore.DB.ExecContext(ctx, query, id, ip)
	if err != nil {
		return fmt.Errorf("failed to insert stat to db: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affect are failed: %w", err)
	}

	if n < 1 {
		return errors.New("no affected rows while saving statistics")
	}
	return nil
}
