package main

import (
	"context"
	"fmt"
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
	db, err := postgres.NewDB(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := storage.NewDB(db)
	a := starter.NewApp(repo)
	h := handler.NewRouter(repo)
	srv := server.NewServer(":8000", h)

	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		a.Serve(ctx, srv)
		wg.Done()
	}()

	select {
	case err := <-srv.Err:
		log.Println(fmt.Errorf("server error: %w", err))
	case <-ctx.Done():
	}
	cancel()
	wg.Wait()
}
