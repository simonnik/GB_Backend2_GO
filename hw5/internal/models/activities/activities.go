package activities

import (
	"database/sql"

	"github.com/simonnik/GB_Backend2_GO/hw5/internal/pool"
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/sharding"
)

type Activity struct {
	Manager *sharding.Manager
	Pool    *pool.Pool
	UserId  int
	Date    string
	Name    string
}

func (a *Activity) connection() (*sql.DB, error) {
	s, err := a.Manager.ShardById(a.UserId)
	if err != nil {
		return nil, err
	}
	return a.Pool.Connection(s.Address)
}

func (a *Activity) Create() error {
	c, err := a.connection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`INSERT INTO "activities" VALUES ($1, $2, $3)`, a.UserId,
		a.Date, a.Name)
	return err
}
func (a *Activity) Read() error {
	c, err := a.connection()
	if err != nil {
		return err
	}
	r := c.QueryRow(`SELECT "name", "date" FROM "activities" WHERE "user_id" =
$1`, a.UserId)
	return r.Scan(
		&a.Name,
		&a.Date,
	)
}
func (a *Activity) Update() error {
	c, err := a.connection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`UPDATE "activities" SET "name" = $2, "date" = $3
WHERE "user_id" = $1`, a.UserId, a.Name, a.Date)
	return err
}
func (a *Activity) Delete() error {
	c, err := a.connection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`DELETE FROM "activities" WHERE "user_id" = $1`, a.UserId)
	return err
}
