package main

import (
	"log/slog"
	"os"
	"url_shortener/internal/config"
	"url_shortener/internal/storage/postgres"

	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	godotenv.Load() // load dotenv file
	cfg := config.MustLoad()

	log := createLogger(cfg.Env)
	log.Info("application has been started")

	storage, err := postgres.New(cfg)
	if err != nil {
		log.Error("fail during loading the storage", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		os.Exit(1)
	}

	_ = storage
}

func createLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
