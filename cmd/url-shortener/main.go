package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/famineBurgund/famiURL/internal/config"
	"github.com/famineBurgund/famiURL/internal/http-server/handlers/url/save"
	"github.com/famineBurgund/famiURL/internal/lib/logger/sl"
	"github.com/famineBurgund/famiURL/internal/storage/postgres"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func main() {
	// TODO: init config: cleanenv
	if err := godotenv.Load("local.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := config.MustLoad()

	// TODO: init logger: log/slog
	log := setupLogger(cfg.Env)

	log.Info("starting url shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enable")

	// TODO: init storage: postgres
	storage, err := postgres.New(cfg.StoragePath)
	if err != nil {
		log.Error("fail init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	// TODO: init router: chi, render
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server error", sl.Err(err))
	}

	log.Info("server stopped")

	// TODO: run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case EnvLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case EnvDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
