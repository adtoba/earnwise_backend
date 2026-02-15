package routes

import (
	"github.com/adtoba/earnwise_backend/src/controllers"
	"github.com/adtoba/earnwise_backend/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type WalletRouteController struct {
	walletController controllers.WalletController
}

func NewWalletRouteController(walletController controllers.WalletController) *WalletRouteController {
	return &WalletRouteController{walletController: walletController}
}

func (rc *WalletRouteController) RegisterWalletRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	router := rg.Group("/wallets")
	router.GET("/", middleware.AuthMiddleware(redisClient), rc.walletController.GetWallet)
}
