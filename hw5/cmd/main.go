package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/simonnik/GB_Backend2_GO/hw5/internal/api"
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/core/activities"
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/core/users"
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/logic/server"
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/pool"
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/repo/datastore"
	activitiesRepo "github.com/simonnik/GB_Backend2_GO/hw5/internal/repo/postgres/activities"
	usersRepo "github.com/simonnik/GB_Backend2_GO/hw5/internal/repo/postgres/users"
	"github.com/simonnik/GB_Backend2_GO/hw5/internal/sharding"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	m := sharding.NewManager(10)
	m.Add(&sharding.Shard{"port=8100 user=test password=test dbname=test sslmode=disable", 0})
	m.Add(&sharding.Shard{"port=8110 user=test password=test dbname=test sslmode=disable", 1})
	m.Add(&sharding.Shard{"port=8120 user=test password=test dbname=test sslmode=disable", 2})

	p := pool.NewPool()
	defer p.Close()

	ds, err := datastore.NewDatastore(m, p)
	if err != nil {
		log.Fatal(err)
	}
	repoUsers := usersRepo.NewRepository(ds)
	us := users.NewService(repoUsers)

	repoActivity := activitiesRepo.NewRepository(ds)
	as := activities.NewService(repoActivity)

	a := api.NewAPI(us, as)
	srv := server.NewAPIServer(a)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	chErr := make(chan error, 1)
	go func() {
		if err := srv.Start(":8080"); err != nil && err != http.ErrServerClosed {
			chErr <- err
		}
		wg.Done()
	}()

	select {
	case err := <-chErr:
		srv.Logger.Error(fmt.Errorf("server error: %w", err))
	case <-ctx.Done():
		srv.Logger.Info("shutdown inited")

		if err := srv.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}
	cancel()
	wg.Wait()
}
