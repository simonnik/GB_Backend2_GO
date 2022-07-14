package main

import "log"

func main() {
	const (
		host = "localhost"
		port = "6379"
	)
	client, err := NewRedisClient(host, port)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	if err := withStructWork(client); err != nil {
		log.Fatal(err)
	}
}
