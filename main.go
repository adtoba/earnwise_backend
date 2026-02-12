package main

import (
	"log"
	"net/http"

	"github.com/adtoba/earnwise_backend.git/src/initializers"
	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine
)

func init() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	initializers.ConnectDB(&config)

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
