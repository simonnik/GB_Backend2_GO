package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/trace"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"

	"github.com/simonnik/GB_Backend2_GO/hw4/app/internal/red"
	"github.com/simonnik/GB_Backend2_GO/hw4/infrastructure/telemetry"
)

var (
	db *sql.DB

	measurable = red.MeasurableHandler

	router = mux.NewRouter()
	web    = http.Server{
		Addr:    ":80",
		Handler: router,
	}
)

type S struct {
	Tr trace.Tracer
}

func main() {
	ctx := context.Background()
	// Настраиваем сборщик трейсов
	tp, err := telemetry.RunTracingCollection(ctx)
	if err != nil {
		panic(fmt.Errorf("error on run tracing collection: %v", err))
	}
	defer func() {
		if err = tp.Shutdown(context.Background()); err != nil {
			panic(fmt.Errorf("failed to stop the traces collector: %v", err))
		}
	}()

	tr := tp.Tracer("server")
	s := &S{
		Tr: tr,
	}

	router.
		HandleFunc("/entities", measurable(s.ListAllEntitiesHandler)).
		Methods(http.MethodGet)
	router.
		HandleFunc("/new-entity", measurable(s.AddEntityHandler)).
		Methods(http.MethodPost)
	router.
		HandleFunc("/get-entity", measurable(s.ReadEntityHandler)).
		Methods(http.MethodGet)
	db, err = sql.Open("mysql", "root:test@tcp(database:3306)/test")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	log.Println("connected to database")
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != http.ErrServerClosed {
			panic(fmt.Errorf("error on listen and serve: %v", err))
		}
	}()

	if err := web.ListenAndServe(); err != http.ErrServerClosed {
		panic(fmt.Errorf("error on listen and serve: %v", err))
	}
}

const sqlInsertEntity = `
  INSERT INTO entities(id, data) VALUES (?, ?)
  `

// AddEntityHandler ...
func (s *S) AddEntityHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.Tr.Start(ctx, "ReadEntityHandler")
	defer span.End()
	res, err := http.Get(fmt.Sprintf("http://acl/identity?token=%s", r.FormValue("token")))
	switch {
	case err != nil:
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	case res.StatusCode != http.StatusOK:
		w.WriteHeader(http.StatusForbidden)
		return
	}
	res.Body.Close()

	// _, err = db.Exec(sqlInsertEntity, r.FormValue("id"), r.FormValue("data"))
	err = red.MeasurableDBExec(db, sqlInsertEntity, r.FormValue("id"), r.FormValue("data"))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

const sqlSelectEntities = `
  SELECT id, data FROM entities;
  `

// ListEntityItemResponse ...
type ListEntityItemResponse struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

// ListAllEntitiesHandler ...
func (s *S) ListAllEntitiesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.Tr.Start(ctx, "ListAllEntitiesHandler")
	defer span.End()
	// rr, err := db.Query(sqlSelectEntities)
	rr, err := red.MeasurableDBQuery(db, sqlSelectEntities)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rr.Close()

	var ii []*ListEntityItemResponse
	for rr.Next() {
		i := &ListEntityItemResponse{}
		err = rr.Scan(&i.ID, &i.Data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ii = append(ii, i)
	}
	bb, err := json.Marshal(ii)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(bb)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

const sqlSelectOneEntity = `
  SELECT id, data FROM entities WHERE id = ?;
  `

// ReadEntityHandler ...
func (s *S) ReadEntityHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.Tr.Start(ctx, "ReadEntityHandler")
	defer span.End()
	// row := db.QueryRow(sqlSelectEntities)
	row := red.MeasurableDBQueryRow(db, sqlSelectOneEntity, r.FormValue("id"))
	i := &ListEntityItemResponse{}
	err := row.Scan(&i.ID, &i.Data)
	if err != nil {
		red.ErrorDBQueryRow()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bb, err := json.Marshal(i)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(bb)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
