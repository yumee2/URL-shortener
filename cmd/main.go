package main

import (
	"url_shortener/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load() // load dotenv file
	cfg := config.MustLoad()
}
