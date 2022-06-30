package starter

import (
	"context"

	"github.com/simonnik/GB_Backend2_GO/internal/logic/storage"
)

// App struct of application
type App struct {
	db *storage.DB
}

// NewApp creates Application instance
func NewApp(store *storage.DB) *App {
	a := &App{
		db: store,
	}
	return a
}

// APIServer is interface to interact with particular server
type APIServer interface {
	Start(db *storage.DB)
	Stop()
}

// Serve the server starts and stops when the context done
func (a *App) Serve(ctx context.Context, hs APIServer) {
	hs.Start(a.db)
	<-ctx.Done()
	hs.Stop()
}
