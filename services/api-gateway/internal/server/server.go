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
	userHandler *handler.UserHandler
	restaurantHandler *handler.RestaurantHandler
	config      apigateway.Env
}

func NewServer(cfg *apigateway.Env, uaClient *client.UAServiceClient, restaurantClient *client.RestaurantServiceClient) *Server {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(GinLogger())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*"},
	}))

	authHandler := handler.NewAuthHandler(uaClient)
	userHandler := handler.NewUserHandler(uaClient)
	restaurantHandler := handler.NewRestaurantHandler(restaurantClient)
	
	return &Server{
		router:      router,
		authHandler: authHandler,
		userHandler: userHandler,
		restaurantHandler: restaurantHandler,
		config:      *cfg,
	}
}

func (s *Server) SetupRoutes() {
	v1 := s.router.Group("/api/v1")

	// Home endpoint
	s.router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the API Gateway, where all your requests find their way!",
		})
	})

	// Health check endpoint
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	// Auth routes
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", s.authHandler.Register)
			auth.POST("/login", s.authHandler.LoginWithEmailAndPassword)
			auth.POST("/google", s.authHandler.LoginWithGoogle)
			auth.POST("/logout", s.authHandler.Logout)
			auth.POST("/refresh", s.authHandler.Refresh)
		}
	}

	// User routes
	{
		user := v1.Group("/user")
		{
			user.GET("/", s.userHandler.GetUser)
			user.GET("/phone-number", s.userHandler.GetPhoneNumber)
			user.POST("/phone-number", s.userHandler.AddPhoneNumber)
			user.PUT("/phone-number", s.userHandler.UpdatePhoneNumber)
			user.DELETE("/phone-number", s.userHandler.RemovePhoneNumber)
			user.GET("/addresses", s.userHandler.GetAddresses)
			user.POST("/addresses", s.userHandler.AddAddress)
			user.DELETE("/addresses", s.userHandler.RemoveAddress)

			user.POST("/be-driver", s.userHandler.BeDriver)
			user.GET("/drivers", s.userHandler.GetDrivers)
			user.DELETE("/driver", s.userHandler.RemoveDriver)
		}
	}

	// Restaurant routes
	{
		restaurant := v1.Group("/restaurants")
		{
			restaurant.POST("/login", s.restaurantHandler.Login)
			restaurant.POST("/register", s.restaurantHandler.RegisterRestaurant)

			restaurant.GET("/", s.restaurantHandler.GetRestaurant)
			restaurant.POST("/", s.restaurantHandler.ListRestaurants)

			restaurant.POST("/menu", s.restaurantHandler.AddMenuItem)
			restaurant.PUT("/menu", s.restaurantHandler.UpdateMenuItem)
			restaurant.DELETE("/menu", s.restaurantHandler.RemoveMenuItem)
		}
	}

	// Order routes
	{
		order := v1.Group("/orders")
		{
			order.POST("/", s.restaurantHandler.PlaceOrder)
			order.GET("/:order_id", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Order details fetched successfully!",
				})
			})
			order.PUT("/:order_id/status", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Order status updated successfully!",
				})
			})
		}
	}

	// Delivery routes
	{
		delivery := v1.Group("/deliveries")
		{
			delivery.POST("/", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Delivery created successfully!",
				})
			})
			delivery.GET("/:delivery_id", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Delivery details fetched successfully!",
				})
			})
			delivery.PUT("/:delivery_id/status", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Delivery status updated successfully!",
				})
			})
		}
	}

	// Payment routes
	{
		payment := v1.Group("/payments")
		{
			payment.POST("/", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Payment processed successfully!",
				})
			})
			payment.GET("/:payment_id", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Payment details fetched successfully!",
				})
			})
		}
	}

}

func (s *Server) Run() error {
	return s.router.Run(":" + s.config.API_GATEWAY_PORT)
}
