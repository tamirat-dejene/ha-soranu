package server

import (
	"github.com/gin-gonic/gin"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/api/handler"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/client"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/config"
)

type Server struct {
	router      *gin.Engine
	authHandler *handler.AuthHandler
	config      config.Config
}

func NewServer(cfg config.Config, authClient *client.AuthClient) *Server {
	router := gin.Default()
	authHandler := handler.NewAuthHandler(authClient)

	return &Server{
		router:      router,
		authHandler: authHandler,
		config:      cfg,
	}
}

func (s *Server) SetupRoutes() {
	v1 := s.router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", s.authHandler.Register)
			auth.POST("/login", s.authHandler.Login)
			auth.POST("/logout", s.authHandler.Logout)
			auth.POST("/refresh", s.authHandler.Refresh)
		}
	}
}

func (s *Server) Run() error {
	return s.router.Run(":" + s.config.ServerPort)
}
