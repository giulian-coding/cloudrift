package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/giulian-coding/cloudrift/internal/containers"
	"github.com/joho/godotenv"
	"github.com/containerd/containerd"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	gin.SetMode(gin.DebugMode)

	router := gin.New()
	router.HandleMethodNotAllowed = true

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error":  "method not allowed",
			"method": c.Request.Method,
		})
	})

	client, err := containerd.New("/run/containerd/containerd.sock", containerd.WithDefaultNamespace("default"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()


	containerRepo := containers.NewMemoryRepository()
	containerSvc := containers.NewService(containerRepo, client)
	containerHandler := containers.NewHandler(containerSvc)

	containerGroup := router.Group("/containers")
	containerHandler.Register(containerGroup)
	

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
