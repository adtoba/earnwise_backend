package main

import (
	"log"
	"net/http"

	"github.com/adtoba/earnwise_backend/src/controllers"
	"github.com/adtoba/earnwise_backend/src/initializers"
	"github.com/adtoba/earnwise_backend/src/migrate"
	"github.com/adtoba/earnwise_backend/src/routes"
	"github.com/adtoba/earnwise_backend/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	server      *gin.Engine
	RedisClient *redis.Client

	AuthController      *controllers.AuthController
	AuthRouteController *routes.AuthRouteController

	CategoryController      *controllers.CategoryController
	CategoryRouteController *routes.CategoryRouteController

	ExpertController      *controllers.ExpertController
	ExpertRouteController *routes.ExpertRouteController

	UserController      *controllers.UserController
	UserRouteController *routes.UserRouteController

	WalletController      *controllers.WalletController
	WalletRouteController *routes.WalletRouteController

	PostController      *controllers.PostController
	PostRouteController *routes.PostRouteController

	ReviewController      *controllers.ReviewController
	ReviewRouteController *routes.ReviewRouteController

	ChatController      *controllers.ChatController
	ChatRouteController *routes.ChatRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	DB := initializers.ConnectDB(&config)
	migrate.Migrate(DB)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Username: config.RedisUsername,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	tokenMaker := utils.NewJWTMaker(config.JWTSecret, RedisClient)

	AuthController = controllers.NewAuthController(DB, tokenMaker, RedisClient)
	AuthRouteController = routes.NewAuthRouteController(*AuthController)

	CategoryController = controllers.NewCategoryController(DB)
	CategoryRouteController = routes.NewCategoryRouteController(*CategoryController)

	WalletController = controllers.NewWalletController(DB)
	WalletRouteController = routes.NewWalletRouteController(*WalletController)

	ExpertController = controllers.NewExpertController(DB, WalletController)
	ExpertRouteController = routes.NewExpertRouteController(*ExpertController)

	UserController = controllers.NewUserController(DB)
	UserRouteController = routes.NewUserRouteController(*UserController)

	PostController = controllers.NewPostController(DB)
	PostRouteController = routes.NewPostRouteController(*PostController)

	ReviewController = controllers.NewReviewController(DB)
	ReviewRouteController = routes.NewReviewRouteController(*ReviewController)

	ChatController = controllers.NewChatController(DB)
	ChatRouteController = routes.NewChatRouteController(*ChatController)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	server.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// server.RedirectTrailingSlash = false

	v1 := server.Group("/api/v1")
	AuthRouteController.RegisterAuthRoutes(v1)
	CategoryRouteController.RegisterCategoryRoutes(v1)
	ExpertRouteController.RegisterExpertRoutes(v1, RedisClient)
	UserRouteController.RegisterUserRoutes(v1, RedisClient)
	WalletRouteController.RegisterWalletRoutes(v1, RedisClient)
	PostRouteController.RegisterPostRoutes(v1, RedisClient)
	ReviewRouteController.RegisterReviewRoutes(v1, RedisClient)
	ChatRouteController.RegisterChatRoutes(v1, RedisClient)

	log.Fatal(server.Run(":" + config.Port))
}
