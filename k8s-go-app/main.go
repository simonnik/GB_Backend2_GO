package main

import (
	"context"
	"fmt"
	"k8s-go-app/config"
	"k8s-go-app/server"
	"k8s-go-app/version"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	launchMode := config.LaunchMode(os.Getenv("LAUNCH_MODE"))
	if len(launchMode) == 0 {
		launchMode = config.LocalEnv
	}
	log.Printf("LAUNCH MODE: %v", launchMode)
	cfg, err := config.Load(launchMode, "./config")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("CONFIG: %+v", cfg)

	info := server.VersionInfo{
		Version: version.Version,
		Commit:  version.Commit,
		Build:   version.Build,
	}
	srv := server.New(info, cfg.Port)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err := srv.Serve(ctx)
		if err != nil {
			log.Println(fmt.Errorf("serve: %w", err))
			return
		}
	}()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-osSignal
	log.Println("OS interrupting signal has receiving")

	cancel()
}
func handler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "Hello, World! Welcome to GeekBrains!\n")
}
