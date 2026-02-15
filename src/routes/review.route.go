package routes

import (
	"github.com/adtoba/earnwise_backend/src/controllers"
	"github.com/adtoba/earnwise_backend/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ReviewRouteController struct {
	reviewController controllers.ReviewController
}

func NewReviewRouteController(reviewController controllers.ReviewController) *ReviewRouteController {
	return &ReviewRouteController{reviewController: reviewController}
}

func (rc *ReviewRouteController) RegisterReviewRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	router := rg.Group("/reviews")
	router.POST("/", middleware.AuthMiddleware(redisClient), rc.reviewController.CreateReview)
	router.GET("/expert/:id", middleware.AuthMiddleware(redisClient), rc.reviewController.GetReviewsByExpertId)
	router.GET("/user/:id", middleware.AuthMiddleware(redisClient), rc.reviewController.GetReviewsByUserId)
}
