package main

import (
	"log"
	"net/http"

	"github.com/adtoba/earnwise_backend/src/initializers"
	"github.com/adtoba/earnwise_backend/src/migrate"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	server      *gin.Engine
	RedisClient *redis.Client
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

	log.Fatal(server.Run(":" + config.Port))
}
