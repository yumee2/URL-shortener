package main

import (
	"log/slog"
	"os"
	"url_shortener/internal/config"
	"url_shortener/internal/http_server/controllers"
	"url_shortener/internal/http_server/routers"
	"url_shortener/internal/services"
	"url_shortener/internal/storage/postgres"

	"github.com/gin-gonic/gin"
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

	r := setupRouter(*storage, log)
	if err := r.Run(cfg.Addres); err != nil {
		log.Error("Failed to start server:", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
	}
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

func setupRouter(storage postgres.Storage, log *slog.Logger) *gin.Engine {
	r := gin.Default()
	urlService := services.NewURLService(storage, log)
	urlController := controllers.NewURLController(urlService, log)

	routers.SetupURLRoutes(r, urlController)
	return r
}
