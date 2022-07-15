package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	const (
		host        = "localhost"
		redisPort   = "6379"
		servicePort = "8080"
	)
	ttl := 1 * time.Hour
	client, err := NewRedisClient(host, redisPort, ttl)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatal(fmt.Errorf("error while closing redis connection: %w", err))
		}
	}()
	//http.HandleFunc("/", client.RootHandler)
	http.HandleFunc("/signup", client.SignupHandler)
	http.HandleFunc("/signup-complete", client.SignupCompleteHandler)
	log.Printf("starting server at :%s", servicePort)
	log.Fatal(http.ListenAndServe(":"+servicePort, nil))
}
