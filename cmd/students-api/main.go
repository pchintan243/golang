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

	"github.com/go-playground/validator/v10"
	"github.com/pchintan243/golang/internal/config"
	"github.com/pchintan243/golang/internal/http/handlers/student"
	"github.com/pchintan243/golang/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage Initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()
	v := validator.New()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteById(storage))
	router.HandleFunc("PUT /api/students", student.Update(storage, v))

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
