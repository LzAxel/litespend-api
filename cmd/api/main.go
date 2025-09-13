package main

import (
	"litespend-api/internal/app"
	"litespend-api/internal/config"
	"litespend-api/internal/pkg/logger"
)

func main() {
	logger.InitLogger()
	cfg := config.GetConfig()
	api := app.NewApp(cfg)

	api.Run()
}
