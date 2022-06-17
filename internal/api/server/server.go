package server

import (
	"context"
	"net/http"
	"time"

	"github.com/simonnik/GB_Backend2_GO/internal/logic/starter"
	"github.com/simonnik/GB_Backend2_GO/internal/logic/storage"
)

var _ starter.APIServer = &Server{}

type Server struct {
	srv http.Server
	db  *storage.DB
}

func NewServer(addr string, h http.Handler) *Server {
	s := &Server{}

	s.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}
	return s
}

func (s *Server) Start(db *storage.DB) {
	s.db = db
	go s.srv.ListenAndServe()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	s.srv.Shutdown(ctx)
	cancel()
}
