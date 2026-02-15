package routes

import (
	"github.com/adtoba/earnwise_backend/src/controllers"
	"github.com/adtoba/earnwise_backend/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type PostRouteController struct {
	postController controllers.PostController
}

func NewPostRouteController(postController controllers.PostController) *PostRouteController {
	return &PostRouteController{postController: postController}
}

func (rc *PostRouteController) RegisterPostRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	router := rg.Group("/posts")
	router.POST("/", middleware.AuthMiddleware(redisClient), rc.postController.CreatePost)
	router.GET("/", middleware.AuthMiddleware(redisClient), rc.postController.GetPosts)
	router.GET("/:id", middleware.AuthMiddleware(redisClient), rc.postController.GetPostById)
	router.GET("/expert/:id", middleware.AuthMiddleware(redisClient), rc.postController.GetPostsByExpertId)
}
