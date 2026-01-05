package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	apigateway "github.com/tamirat-dejene/ha-soranu/services/api-gateway"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/api/handler"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/client"
)

type Server struct {
	router              *gin.Engine
	authHandler         *handler.AuthHandler
	userHandler         *handler.UserHandler
	restaurantHandler   *handler.RestaurantHandler
	notificationHandler *handler.NotificationHandler
	config              apigateway.Env
}

func NewServer(cfg *apigateway.Env, uaClient *client.UAServiceClient, restaurantClient *client.RestaurantServiceClient, notificationClient *client.NotificationServiceClient) *Server {
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
	notificationHandler := handler.NewNotificationHandler(notificationClient)

	return &Server{
		router:              router,
		authHandler:         authHandler,
		userHandler:         userHandler,
		restaurantHandler:   restaurantHandler,
		notificationHandler: notificationHandler,
		config:              *cfg,
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
			user.POST("/drivers", s.userHandler.GetDrivers)
			user.DELETE("/drivers", s.userHandler.RemoveDriver)

			// User Notifications
			user.GET("/:user_id/notifications", s.notificationHandler.GetUserNotifications)
		}
	}

	// Restaurant routes
	{
		restaurant := v1.Group("/restaurants", AuthMiddleware(&s.config))
		{
			restaurant.POST("/login", s.restaurantHandler.Login)
			restaurant.POST("/register", s.restaurantHandler.RegisterRestaurant)

			restaurant.GET("/", s.restaurantHandler.GetRestaurant)
			restaurant.POST("/", s.restaurantHandler.ListRestaurants)

			// Menu routes for restaurants
			restaurant.POST("/menu", s.restaurantHandler.AddMenuItem)
			restaurant.PUT("/menu", s.restaurantHandler.UpdateMenuItem)
			restaurant.DELETE("/menu", s.restaurantHandler.RemoveMenuItem)

			// Order routes for restaurants
			restaurant.GET("/orders", s.restaurantHandler.GetOrders)
			restaurant.POST("/orders", s.restaurantHandler.PlaceOrder)
			restaurant.PUT("/:restaurant_id/orders/:order_id/status", s.restaurantHandler.UpdateOrderStatus)
			restaurant.PUT("orders/:order_id/ship", s.restaurantHandler.ShipOrder)
			restaurant.GET("/orders/:order_id", s.restaurantHandler.GetOrder)

			// Restaurant Notifications
			restaurant.GET("/:restaurant_id/notifications", s.notificationHandler.GetRestaurantNotifications)
		}
	}

	// Notification routes
	{
		notification := v1.Group("/notifications", AuthMiddleware(&s.config))
		{
			notification.PUT("/:notification_id/read", s.notificationHandler.MarkAsRead)
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
		payment := v1.Group("/payments", AuthMiddleware(&s.config))
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
