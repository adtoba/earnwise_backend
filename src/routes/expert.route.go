package routes

import (
	"github.com/adtoba/earnwise_backend/src/controllers"
	"github.com/adtoba/earnwise_backend/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ExpertRouteController struct {
	expertController controllers.ExpertController
}

func NewExpertRouteController(expertController controllers.ExpertController) *ExpertRouteController {
	return &ExpertRouteController{expertController: expertController}
}

func (rc *ExpertRouteController) RegisterExpertRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	router := rg.Group("/experts")
	router.POST("/", middleware.AuthMiddleware(redisClient), rc.expertController.CreateExpertProfile)
	router.GET("/dashboard", middleware.AuthMiddleware(redisClient), rc.expertController.GetExpertDashboard)
	router.GET("/:id", middleware.AuthMiddleware(redisClient), rc.expertController.GetExpertProfileById)
	router.GET("/", middleware.AuthMiddleware(redisClient), rc.expertController.GetExpertProfile)
	router.PUT("/rate", middleware.AuthMiddleware(redisClient), rc.expertController.UpdateExpertRate)
	router.PUT("/socials", middleware.AuthMiddleware(redisClient), rc.expertController.UpdateExpertSocials)
	router.PUT("/details", middleware.AuthMiddleware(redisClient), rc.expertController.UpdateExpertDetails)
	router.PUT("/availability", middleware.AuthMiddleware(redisClient), rc.expertController.UpdateExpertAvailability)
	router.GET("/category/:category", middleware.AuthMiddleware(redisClient), rc.expertController.GetExpertsByCategory)
	router.GET("/recommended", middleware.AuthMiddleware(redisClient), rc.expertController.GetRecommendedTopExperts)
}
