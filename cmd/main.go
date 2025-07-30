package main

import (
	"log/slog"
	"os"
	"url_shortener/internal/config"

	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	godotenv.Load() // load dotenv file
	cfg := config.MustLoad()

	logger := createLogger(cfg.Env)
	logger.Info("application has been started")
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
