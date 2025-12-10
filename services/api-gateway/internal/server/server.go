package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	apigateway "github.com/tamirat-dejene/ha-soranu/services/api-gateway"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/api/handler"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/client"
)

type Server struct {
	router      *gin.Engine
	authHandler *handler.AuthHandler
	config      apigateway.Env
}

func NewServer(cfg *apigateway.Env, authClient *client.AuthClient) *Server {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(GinLogger())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*"},
	}))

	authHandler := handler.NewAuthHandler(authClient)
	return &Server{
		router:      router,
		authHandler: authHandler,
		config:      *cfg,
	}
}

func (s *Server) SetupRoutes() {
	v1 := s.router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", s.authHandler.Register)
			auth.POST("/login", s.authHandler.Login)
			auth.POST("/google", s.authHandler.LoginWithGoogle)
			auth.POST("/logout", s.authHandler.Logout)
			auth.POST("/refresh", s.authHandler.Refresh)
		}
	}

	// Health check endpoint
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	// Home endpoint
	s.router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the API Gateway, where all your requests find their way!",
		})
	})
}

func (s *Server) Run() error {
	return s.router.Run(":" + s.config.API_GATEWAY_PORT)
}
