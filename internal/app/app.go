package app

import (
	"context"
	"litespend-api/internal/config"
	"litespend-api/internal/httpsrv"
	"litespend-api/internal/repository"
	"litespend-api/internal/repository/databases"
	"litespend-api/internal/service"
	"litespend-api/internal/session"
	"log/slog"
)

type App struct {
	config     config.Config
	repository *repository.Repository
	service    *service.Service
	server     *httpsrv.Server
}

func NewApp(cfg config.Config) *App {
	ctx := context.Background()
	slog.InfoContext(ctx, "Starting app")

	psql, err := databases.GetPostgresDB(ctx, cfg.Postgres)
	if err != nil {
		panic(err)
	}

	repo := repository.NewRepository(psql)
	services := service.NewService(repo)

	sessionManager := session.NewSessionManager(session.NewSessionInMemoryStore())

	server := httpsrv.NewServer(cfg.Server, services, sessionManager)

	app := App{
		config:     cfg,
		repository: repo,
		service:    services,
		server:     server,
	}

	return &app
}

func (a *App) Run() {
	slog.InfoContext(context.Background(), "Starting server")
	a.server.Run()
}
