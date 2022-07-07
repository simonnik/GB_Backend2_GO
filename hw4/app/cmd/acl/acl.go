package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/simonnik/GB_Backend2_GO/hw4/app/internal/red"
)

var (
	measurable = red.MeasurableHandler

	router = mux.NewRouter()
	web    = http.Server{
		Handler: router,
	}
)

func init() {
	router.
		HandleFunc("/identity", measurable(GetIdentityHandler)).
		Methods(http.MethodGet)
}

func main() {
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

// GetIdentityHandler ...
func GetIdentityHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("token") == "admin_secret_token" {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}
