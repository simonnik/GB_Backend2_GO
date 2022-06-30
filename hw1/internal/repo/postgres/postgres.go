package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/simonnik/GB_Backend2_GO/internal/entity"
	"github.com/simonnik/GB_Backend2_GO/internal/logic/storage"
)

// PGRepo works with PG
type PGRepo struct {
	db *pgxpool.Pool
}

// Checking if the interface matches
var _ storage.Repo = &PGRepo{}

// NewDB connects to DB and produces PGRepo
func NewDB(ctx context.Context, connStr string) (*PGRepo, error) {
	// pool, err := getConn(connStr)
	config, err := getPGXPoolConfig(connStr)
	conn, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &PGRepo{
		db: conn,
	}, nil
}

// Close connection to DB
func (pg *PGRepo) Close() {
	pg.db.Close()
}

// CreateUser executes INSERT query into USERS table
func (pg *PGRepo) CreateUser(ctx context.Context, u entity.User) error {
	_, err := pg.db.Exec(ctx, "INSERT INTO users(uuid, user_name, email) VALUES ($1, $2, $3)", u.ID, u.Name, u.Email)
	return err
}

// CreateGroup executes INSERT query into GROUPS table
func (pg *PGRepo) CreateGroup(ctx context.Context, g entity.Group) error {
	_, err := pg.db.Exec(ctx, "INSERT INTO groups(uuid, group_name, group_type, descr) VALUES ($1, $2, $3, $4)", g.ID, g.Name, g.Type, g.Description)
	return err
}

// AddToGroup adds a row (gid, uid) into membership table
func (pg *PGRepo) AddToGroup(ctx context.Context, uid, gid uuid.UUID) error {
	_, err := pg.db.Exec(ctx, "INSERT INTO membership(group_uuid, user_uuid) VALUES ($1, $2)", gid, uid)
	return err
}

// RemoveFromGroup removes a row (gid, uid) from membership table
func (pg *PGRepo) RemoveFromGroup(ctx context.Context, uid, gid uuid.UUID) error {
	_, err := pg.db.Exec(ctx, "DELETE FROM membership WHERE group_uuid=$1 AND user_uuid=$2)", gid, uid)
	return err
}

// SearchUser search user by name and by member in groups if set
func (pg *PGRepo) SearchUser(ctx context.Context, name string, gids []uuid.UUID) ([]entity.User, error) {

	name += "%"
	var rows pgx.Rows
	var err error

	switch len(gids) {
	case 0:
		// search by name only
		query := "SELECT * FROM users WHERE user_name LIKE $1"
		rows, err = pg.db.Query(ctx, query, name)
	case 1:
		// search by name and one guid
		query := `
SELECT uuid, user_name, email FROM users WHERE uuid IN (SELECT u.uuid FROM users u WHERE u.user_name LIKE $1
INTERSECT
SELECT m.user_uuid FROM membership m WHERE m.group_uuid=$2);
`
		rows, err = pg.db.Query(ctx, query, name, gids[0])

	case 2:
		// search by name and two guids
		query := `
SELECT uuid, user_name, email FROM users WHERE uuid IN (SELECT u.uuid FROM users u WHERE u.user_name LIKE $1
INTERSECT
SELECT m1.user_uuid FROM membership m1 WHERE m1.group_uuid=$2
INTERSECT
SELECT m2.user_uuid FROM membership m2 WHERE m2.group_uuid=$3);
`
		rows, err = pg.db.Query(ctx, query, name, gids[0], gids[1])

	case 3:
		// search by name and thee guids
		query := `
SELECT uuid, user_name, email FROM users WHERE uuid IN (SELECT u.uuid FROM users u WHERE u.user_name LIKE $1
INTERSECT
SELECT m1.user_uuid FROM membership m1 WHERE m1.group_uuid=$2
INTERSECT
SELECT m2.user_uuid FROM membership m2 WHERE m2.group_uuid=$3
INTERSECT
SELECT m3.user_uuid FROM membership m3 WHERE m3.group_uuid=$4);
`
		rows, err = pg.db.Query(ctx, query, name, gids[0], gids[1], gids[2])
	default:
		return nil, errors.New("too many guids passed, 3 is maximum")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]entity.User, 0)
	for rows.Next() {
		u := &entity.User{}
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, fmt.Errorf("failed to get rows on given phrase: %w", err)
		}
		users = append(users, *u)
	}
	return users, nil

}

// SearchGroup search group by name and by members if set
func (pg *PGRepo) SearchGroup(ctx context.Context, name string, uids []uuid.UUID) ([]entity.Group, error) {

	name += "%"
	var rows pgx.Rows
	var err error

	switch len(uids) {
	case 0:
		// search by name only
		query := "SELECT uuid, group_name, group_type, descr FROM groups WHERE group_name LIKE $1"
		rows, err = pg.db.Query(ctx, query, name)
	case 1:
		// search by name and one guid
		query := `
SELECT uuid, group_name, group_type, descr FROM users WHERE uuid IN (SELECT u.uuid FROM users u WHERE u.user_name LIKE $1
INTERSECT
SELECT m.group_uuid FROM membership m WHERE m.user_uuid=$2);
`
		rows, err = pg.db.Query(ctx, query, name, uids[0])

	case 2:
		// search by name and two guids
		query := `
SELECT uuid, group_name, group_type, descr FROM users WHERE uuid IN (SELECT u.uuid FROM users u WHERE u.user_name LIKE $1
INTERSECT
SELECT m1.group_uuid FROM membership m1 WHERE m1.user_uuid=$2
INTERSECT
SELECT m2.group_uuid FROM membership m2 WHERE m2.user_uuid=$3);
`
		rows, err = pg.db.Query(ctx, query, name, uids[0], uids[1])

	case 3:
		// search by name and thee guids
		query := `
SELECT uuid, group_name, group_type, descr FROM groups WHERE uuid IN (SELECT g.uuid FROM groups g WHERE g.group_name LIKE $1
INTERSECT
SELECT m1.group_uuid FROM membership m1 WHERE m1.user_uuid=$2
INTERSECT
SELECT m2.group_uuid FROM membership m2 WHERE m2.user_uuid=$3
INTERSECT
SELECT m3.group_uuid FROM membership m3 WHERE m3.user_uuid=$4);
`
		rows, err = pg.db.Query(ctx, query, name, uids[0], uids[1], uids[2])
	default:
		// too many guis passed
		return nil, errors.New("too many guids passed, 3 is maximum")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]entity.Group, 0)
	for rows.Next() {
		g := &entity.Group{}
		if err := rows.Scan(&g.ID, &g.Name, &g.Type, &g.Description); err != nil {
			return nil, fmt.Errorf("failed to get rows on given phrase: %w", err)
		}
		groups = append(groups, *g)
	}
	return groups, nil
}

func getPGXPoolConfig(connStr string) (*pgxpool.Config, error) {

	// "postgresql://user:password@host:port/dbname"

	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create the PGX pool config from connection string: %w", err)
	}
	cfg.ConnConfig.ConnectTimeout = time.Second * 3
	return cfg, nil
}
