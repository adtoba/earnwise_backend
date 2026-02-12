package routes

import (
	"github.com/adtoba/earnwise_backend/src/controllers"
	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) *AuthRouteController {
	return &AuthRouteController{authController: authController}
}

func (rc *AuthRouteController) RegisterAuthRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/auth")
	router.POST("/login", rc.authController.Login)
	router.POST("/register", rc.authController.Register)
}
