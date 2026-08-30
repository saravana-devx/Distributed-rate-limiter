package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ratelimiter/internal/bootstrap"
	"ratelimiter/internal/config"
	"syscall"
	"time"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatalf("No .env file found or having an issue with env variables : %v", err)
	}
	cfg := config.Get()
	app, err := bootstrap.New(cfg)
	if err != nil {
		log.Fatalf("Bootstrap failed : %v", err)
	}
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: app.Router(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start : %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server stop...")
	}
	app.Stop()
}
