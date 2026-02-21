package routes

import (
	"github.com/adtoba/earnwise_backend/src/controllers"
	"github.com/adtoba/earnwise_backend/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ChatRouteController struct {
	chatController controllers.ChatController
}

func NewChatRouteController(chatController controllers.ChatController) *ChatRouteController {
	return &ChatRouteController{chatController: chatController}
}

func (rc *ChatRouteController) RegisterChatRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	router := rg.Group("/chats")
	router.POST("/", middleware.AuthMiddleware(redisClient), rc.chatController.CreateChat)
	router.GET("/:id/messages", middleware.AuthMiddleware(redisClient), rc.chatController.GetChatMessages)
	router.POST("/:id/messages", middleware.AuthMiddleware(redisClient), rc.chatController.CreateMessage)
	router.GET("/user", middleware.AuthMiddleware(redisClient), rc.chatController.GetUserChats)
	router.GET("/expert", middleware.AuthMiddleware(redisClient), rc.chatController.GetExpertChats)
}
