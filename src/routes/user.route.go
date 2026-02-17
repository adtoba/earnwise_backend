package routes

import (
	"github.com/adtoba/earnwise_backend/src/controllers"
	"github.com/adtoba/earnwise_backend/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewUserRouteController(userController controllers.UserController) *UserRouteController {
	return &UserRouteController{userController: userController}
}

func (rc *UserRouteController) RegisterUserRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	router := rg.Group("/users")
	router.GET("/:id", rc.userController.GetUserById)
	router.PUT("/", middleware.AuthMiddleware(redisClient), rc.userController.UpdateUser)
	router.GET("/profile", middleware.AuthMiddleware(redisClient), rc.userController.GetUserProfile)
	router.POST("/save-expert", middleware.AuthMiddleware(redisClient), rc.userController.SaveExpert)
	router.GET("/saved-experts", middleware.AuthMiddleware(redisClient), rc.userController.GetSavedExperts)
}
