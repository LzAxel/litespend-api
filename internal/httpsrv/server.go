package httpsrv

import (
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"litespend-api/internal/config"
	"litespend-api/internal/httpsrv/router"
	"litespend-api/internal/service"
	"litespend-api/internal/session"
	"log/slog"
	"net"
)

type Server struct {
	gin            *gin.Engine
	config         config.ServerConfig
	service        *service.Service
	router         *router.Router
	sessionManager *session.SessionManager
}

func NewServer(config config.ServerConfig, services *service.Service, sessionManager *session.SessionManager) *Server {
	server := &Server{
		gin:            gin.Default(),
		config:         config,
		service:        services,
		router:         router.NewRouter(services),
		sessionManager: sessionManager,
	}

	server.setup()

	return server
}

func (s *Server) setup() {
	s.gin.Use(
		gin.Recovery(),
		sloggin.New(slog.Default()),
		s.sessionManager.LoadAndSave,
	)

	apiv1 := s.gin.Group("/api/v1")
	{
		auth := apiv1.Group("/user")
		{
			auth.POST("/register", s.router.User.Register)
		}
	}
}

func (s *Server) Run() error {
	return s.gin.Run(net.JoinHostPort(s.config.Host, s.config.Port))
}
