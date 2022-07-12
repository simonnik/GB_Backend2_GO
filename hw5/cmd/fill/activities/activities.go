package main

import (
	"fmt"

	_ "github.com/lib/pq"

	"github.com/simonnik/GB_Backend2_GO/hw5/internal/models/activities"
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/pool"
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/sharding"
)

func main() {
	m := sharding.NewManager(10)
	p := pool.NewPool()
	m.Add(&sharding.Shard{"port=8100 user=test password=test dbname=test sslmode=disable", 0})
	m.Add(&sharding.Shard{"port=8110 user=test password=test dbname=test sslmode=disable", 1})
	m.Add(&sharding.Shard{"port=8120 user=test password=test dbname=test sslmode=disable", 2})
	aa := []*activities.Activity{
		{m, p, 1, "2020-12-12 06:10:10", "Registered"},
		{m, p, 10, "2020-12-10 15:43:10", "Registered"},
		{m, p, 13, "2020-11-10 21:09:56", "Registered"},
		{m, p, 25, "2020-11-01 15:43:10", "Registered"},
	}
	for _, a := range aa {
		err := a.Create()
		if err != nil {
			fmt.Println(fmt.Errorf("error on create activity %v: %w", a, err))
		}
	}
}
