package httpsrv

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"litespend-api/internal/config"
	"litespend-api/internal/httpsrv/middleware"
	"litespend-api/internal/httpsrv/router"
	"litespend-api/internal/repository"
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
	repository     *repository.Repository
}

func NewServer(config config.ServerConfig, services *service.Service, sessionManager *session.SessionManager, repo *repository.Repository) *Server {
	server := &Server{
		gin:            gin.Default(),
		config:         config,
		service:        services,
		router:         router.NewRouter(services, sessionManager),
		sessionManager: sessionManager,
		repository:     repo,
	}

	server.setup()

	return server
}

func (s *Server) setup() {
	s.gin.Use(
		gin.Recovery(),
		sloggin.New(slog.Default()),
		cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:5173"},
			AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "PATCH", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			AllowFiles:       true,
			AllowCredentials: true,
		}),
		s.sessionManager.LoadAndSave,
	)

	apiv1 := s.gin.Group("/api/v1")
	{
		auth := apiv1.Group("/user")
		{
			auth.POST("/register", s.router.User.Register)
			auth.POST("/login", s.router.User.Login)
			auth.POST("/logout", s.router.User.Logout)
		}

		transactions := apiv1.Group("/transactions")
		transactions.Use(middleware.RequireAuth(s.sessionManager, s.repository.UserRepository))
		{
			transactions.POST("", s.router.Transaction.CreateTransaction)
			transactions.GET("", s.router.Transaction.GetTransactions)
			transactions.GET("/statistics/balance", s.router.Transaction.GetBalanceStatistics)
			transactions.GET("/statistics/categories", s.router.Transaction.GetCategoryStatistics)
			transactions.GET("/statistics/periods", s.router.Transaction.GetPeriodStatistics)
			transactions.GET("/:id", s.router.Transaction.GetTransaction)
			transactions.PUT("/:id", s.router.Transaction.UpdateTransaction)
			transactions.DELETE("/:id", s.router.Transaction.DeleteTransaction)
		}

		categories := apiv1.Group("/categories")
		categories.Use(middleware.RequireAuth(s.sessionManager, s.repository.UserRepository))
		{
			categories.POST("", s.router.Category.CreateCategory)
			categories.GET("", s.router.Category.GetCategories)
			categories.GET("/:id", s.router.Category.GetCategory)
			categories.PUT("/:id", s.router.Category.UpdateCategory)
			categories.DELETE("/:id", s.router.Category.DeleteCategory)
		}

		budgets := apiv1.Group("/budgets")
		budgets.Use(middleware.RequireAuth(s.sessionManager, s.repository.UserRepository))
		{
			budgets.POST("", s.router.Budget.CreateBudget)
			budgets.GET("", s.router.Budget.GetBudgets)
			budgets.GET("/period", s.router.Budget.GetBudgetsByPeriod)
			budgets.GET("/:id", s.router.Budget.GetBudget)
			budgets.PUT("/:id", s.router.Budget.UpdateBudget)
			budgets.DELETE("/:id", s.router.Budget.DeleteBudget)
		}

		imports := apiv1.Group("/import")
		imports.Use(middleware.RequireAuth(s.sessionManager, s.repository.UserRepository))
		{
			imports.POST("/parse", s.router.Import.ParseExcelFile)
			imports.POST("/data", s.router.Import.ImportData)
		}

		admin := apiv1.Group("/admin")
		admin.Use(middleware.RequireAuth(s.sessionManager, s.repository.UserRepository))
		admin.Use(middleware.RequireAdmin())
		{
			admin.POST("/sessions/revoke", s.router.Auth.RevokeSession)
			admin.GET("/sessions/info", s.router.Auth.GetSessionInfo)
		}
	}
}

func (s *Server) Run() error {
	return s.gin.Run(net.JoinHostPort(s.config.Host, s.config.Port))
}
