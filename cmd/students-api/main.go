package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/pchintan243/golang/internal"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	// Key: "addr", Value: cfg.Addr
	slog.Info("Server Started", slog.String("address", cfg.Addr))
	fmt.Printf("Server Started: %s", cfg.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start Server")
		}
	}()

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shut down!", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}
