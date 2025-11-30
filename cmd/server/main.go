package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vladopadikk/go-chat/internal/auth"
	"github.com/vladopadikk/go-chat/internal/config"
	"github.com/vladopadikk/go-chat/internal/database"
	"github.com/vladopadikk/go-chat/internal/user"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)
	defer db.Close()

	router := gin.Default()

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg)
	authHandler := auth.NewHandler(authService)

	api := router.Group("/api")

	user.RegisterRoutes(api, userHandler)
	auth.RegisterRoutes(api, authHandler)

	protected := api.Group("")
	protected.Use(auth.AuthMiddleware(cfg.JWTSecret))

	router.Run(":" + cfg.AppPort)

}
