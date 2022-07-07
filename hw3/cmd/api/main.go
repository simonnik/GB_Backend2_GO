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

	"github.com/simonnik/GB_Backend2_GO/hw3/internal/api"
	"github.com/simonnik/GB_Backend2_GO/hw3/internal/config"
	"github.com/simonnik/GB_Backend2_GO/hw3/internal/core/check"
	"github.com/simonnik/GB_Backend2_GO/hw3/internal/core/links"
	"github.com/simonnik/GB_Backend2_GO/hw3/internal/logic/server"
	"github.com/simonnik/GB_Backend2_GO/hw3/internal/repo/datastore"
	checkRepo "github.com/simonnik/GB_Backend2_GO/hw3/internal/repo/postgres/check"
	linksRepo "github.com/simonnik/GB_Backend2_GO/hw3/internal/repo/postgres/links"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	cfg, err := config.BuildConfig()
	if err != nil {
		log.Fatal(err)
	}

	ds, err := datastore.NewDatastore(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := ds.Close(); err != nil {
			log.Println(fmt.Errorf("closing db conn: %w", err))
		}
	}()

	token := links.NewToken(cfg.HashMinLength, cfg.HashSalt)
	repoLinks := linksRepo.NewRepository(ds)
	ls := links.NewService(repoLinks, token, cfg.JWT.Secret)

	repoCheck := checkRepo.NewRepository(ds)
	ch := check.NewService(repoCheck)

	a := api.NewAPI(cfg.Host, ls, ch)
	srv := server.NewAPIServer(cfg, a)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	chErr := make(chan error, 1)
	go func() {
		if err := srv.Start(fmt.Sprintf(":%d", cfg.Port)); err != nil && err != http.ErrServerClosed {
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
