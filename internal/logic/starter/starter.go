package starter

import (
	"context"

	"sync"

	"github.com/simonnik/GB_Backend2_GO/internal/logic/storage"
)

type App struct {
	db *storage.DB
}

func NewApp(stor *storage.DB) *App {
	a := &App{
		db: stor,
	}
	return a
}

type APIServer interface {
	Start(db *storage.DB)
	Stop()
}

func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup, hs APIServer) {
	defer wg.Done()
	hs.Start(a.db)
	<-ctx.Done()
	hs.Stop()
}
