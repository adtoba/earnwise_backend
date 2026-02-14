package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/adtoba/earnwise_backend/src/initializers"
	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func IsTokenBlacklisted(tokenString string, redisClient *redis.Client) bool {
	val, err := redisClient.Get(context.Background(), "blacklist:"+tokenString).Result()
	return err == nil && val == "revoked"
}

func AuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var accessToken string
		config, err := initializers.LoadConfig(".")

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Internal Server Error", nil))
			return
		}

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", nil))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		}

		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", nil))
			return
		}

		if IsTokenBlacklisted(accessToken, redisClient) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", nil))
			return
		}

		payload, err := jwt.ParseWithClaims(accessToken, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("invalid token signing method")
			}
			return []byte(config.JWT_SECRET), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", nil))
			return
		}

		claims, ok := payload.Claims.(*models.UserClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", nil))
			return
		}

		c.Set("user_id", claims.ID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}
