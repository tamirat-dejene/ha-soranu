package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	apigateway "github.com/tamirat-dejene/ha-soranu/services/api-gateway"
	jwtvalidator "github.com/tamirat-dejene/ha-soranu/shared/pkg/auth/jwtvalidator"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"go.uber.org/zap"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			logger.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Duration("latency", latency),
			)
		}
	}
}

// AuthMiddleware verifies access tokens using the shared RSA validator and injects claims into the Gin context.
func AuthMiddleware(cfg *apigateway.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		if cfg.ACCESS_TOKEN_PUBLIC_KEY == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token verifier not configured"})
			return
		}

		claims, err := jwtvalidator.ValidateAccessToken(cfg.ACCESS_TOKEN_PUBLIC_KEY, parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Make claims available to downstream handlers.
		c.Set("user_email", claims.UserEmail)
		c.Set("claims", claims)

		c.Next()
	}
}
