package main

import (
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/giulian-coding/cloudrift/internal/auth"
	"github.com/giulian-coding/cloudrift/internal/containers"
	"github.com/joho/godotenv"

	"github.com/gin-contrib/sessions/cookie"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	gin.SetMode(gin.DebugMode)

	store := cookie.NewStore([]byte("secret"))
	router := gin.New()
	router.Use(sessions.Sessions("test", store))

	privateGroup := router.Group("/private")
	err = auth.UseKindeAuth(
		privateGroup,
		os.Getenv("KINDE_DOMAIN"),
		os.Getenv("KINDE_CLIENT_ID"),
		os.Getenv("KINDE_CLIENT_SECRET"),
		os.Getenv("KINDE_REDIRECT_URI"),
	)
	if err != nil {
		log.Fatal("Error loading KindeAuth")
	}

	privateGroup.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Du bist angemeldet.")
	})

	containerRepo := containers.NewMemoryRepository()
	containerSvc := containers.NewService(containerRepo)
	containerHandler := containers.NewHandler(containerSvc)

	containerGroup := router.Group("/containers")
	containerHandler.Register(containerGroup)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
