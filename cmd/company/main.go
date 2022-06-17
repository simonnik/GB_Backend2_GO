package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/simonnik/GB_Backend2_GO/internal/api/handler"
	"github.com/simonnik/GB_Backend2_GO/internal/api/server"
	"github.com/simonnik/GB_Backend2_GO/internal/logic/starter"
	"github.com/simonnik/GB_Backend2_GO/internal/logic/storage"
	"github.com/simonnik/GB_Backend2_GO/internal/repo/postgres"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	// "postgresql://user:password@host:port/dbname"
	connString := `postgresql://appuser:appuser@127.0.0.1:5432/app_db`
	DB, err := postgres.NewDB(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}

	repo := storage.NewDB(DB)
	a := starter.NewApp(repo)
	h := handler.NewRouter(repo)
	srv := server.NewServer(":8000", h)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go a.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
