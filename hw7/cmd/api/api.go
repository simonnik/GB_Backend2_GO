package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/mediocregopher/radix/v3"
	"gocloud.dev/pubsub"
	//_ "gocloud.dev/pubsub/kafkapubsub"
	_ "gocloud.dev/pubsub/rabbitpubsub"
)

var (
	router = mux.NewRouter()
	web    = http.Server{
		Addr:    ":80",
		Handler: router,
	}
	t        *pubsub.Topic
	s        *radix.Pool
	timeout  = 10 * time.Second
	connFunc = func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr,
			radix.DialTimeout(timeout),
		)
	}
)

func init() {
	router.
		HandleFunc("/rate", PostRateHandler).
		Methods(http.MethodPost)
	router.
		HandleFunc("/total", GetTotalHandler).
		Methods(http.MethodGet)
}
func main() {
	if err := web.ListenAndServe(); err != http.ErrServerClosed {
		panic(fmt.Errorf("error on listen and serve: %v", err))
	}
}
func GetTotalHandler(w http.ResponseWriter, r *http.Request) {
	var rates []string
	err := storage().Do(radix.Cmd(&rates, "LRANGE", "result", "0", "10"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(rates) == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	var sum int
	for _, rate := range rates {
		v, err := strconv.Atoi(rate)
		if err != nil {
			continue
		}
		sum += v
	}
	result := float64(sum) / float64(len(rates))
	_, _ = w.Write([]byte(fmt.Sprintf("%.2f", result)))
}
func PostRateHandler(w http.ResponseWriter, r *http.Request) {
	rate := r.FormValue("rate")
	if _, err := strconv.Atoi(rate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	err := topic(ctx).Send(ctx, &pubsub.Message{
		Body: []byte(rate),
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func topic(ctx context.Context) *pubsub.Topic {
	if t != nil {
		return t
	}
	var err error
	//t, err = pubsub.OpenTopic(ctx, "kafka://rates")
	t, err = pubsub.OpenTopic(ctx, "rabbit://rates")
	if err != nil {
		panic(err)
	}
	return t
}
func storage() *radix.Pool {
	if s != nil {
		return s
	}
	var err error
	s, err = radix.NewPool("tcp", "rmq:6379", 1, radix.PoolConnFunc(connFunc))
	if err != nil {
		panic(err)
	}
	return s
}
