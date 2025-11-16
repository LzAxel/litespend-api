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

	psqlPool, err := databases.GetPostgresPool(ctx, cfg.Postgres)
	if err != nil {
		panic(err)
	}
	psql := databases.GetPostgresDB(psqlPool)

	repo := repository.NewRepository(psql)
	sessionManager := session.NewSessionManager(session.NewSessionPostgresStore(psqlPool))
	services := service.NewService(repo, sessionManager)

	server := httpsrv.NewServer(cfg.Server, services, sessionManager, repo)

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
