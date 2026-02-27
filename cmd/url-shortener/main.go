package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/famineBurgund/famiURL/internal/config"
	"github.com/famineBurgund/famiURL/internal/lib/logger/sl"
	"github.com/famineBurgund/famiURL/internal/storage/postgres"
	"github.com/joho/godotenv"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func main() {
	// TODO: init config: cleanenv
	if err := godotenv.Load("C:/Users/fami/golang/famiURL/local.env"); err != nil {
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
