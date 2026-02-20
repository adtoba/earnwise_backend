package middleware

import (
	"net/http"

	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/gin-gonic/gin"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.MustGet("user_role")
		if userRole != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse("You are not authorized to access this resource", nil))
			return
		}
		c.Next()
	}
}
