package routes

import (
	"github.com/adtoba/earnwise_backend/src/controllers"
	"github.com/adtoba/earnwise_backend/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type CallRouteController struct {
	callController controllers.CallController
}

func NewCallRouteController(callController controllers.CallController) *CallRouteController {
	return &CallRouteController{callController: callController}
}

func (rc *CallRouteController) RegisterCallRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	router := rg.Group("/calls")
	router.POST("/", middleware.AuthMiddleware(redisClient), rc.callController.CreateCall)
	router.GET("/user", middleware.AuthMiddleware(redisClient), rc.callController.GetUserCalls)
	router.GET("/expert", middleware.AuthMiddleware(redisClient), rc.callController.GetExpertCalls)
	router.PUT("/:id/accept", middleware.AuthMiddleware(redisClient), rc.callController.AcceptCall)
}
