package main

import (
	"fmt"

	_ "github.com/lib/pq"

	"github.com/simonnik/GB_Backend2_GO/hw5/internal/models/user"
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/pool"
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/sharding"
)

func main() {
	m := sharding.NewManager(10)
	p := pool.NewPool()
	m.Add(&sharding.Shard{"port=8100 user=test password=test dbname=test sslmode=disable", 0})
	m.Add(&sharding.Shard{"port=8110 user=test password=test dbname=test sslmode=disable", 1})
	m.Add(&sharding.Shard{"port=8120 user=test password=test dbname=test sslmode=disable", 2})
	uu := []*user.User{
		{m, p, 1, "Joe Biden", 78, 10},
		{m, p, 10, "Jill Biden", 69, 1},
		{m, p, 13, "Donald Trump", 74, 25},
		{m, p, 25, "Melania Trump", 78, 13},
	}
	for _, u := range uu {
		err := u.Create()
		if err != nil {
			fmt.Println(fmt.Errorf("error on create user %v: %w", u, err))
		}
	}
}
