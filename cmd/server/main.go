package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vladopadikk/go-chat/internal/auth"
	"github.com/vladopadikk/go-chat/internal/chat"
	"github.com/vladopadikk/go-chat/internal/config"
	"github.com/vladopadikk/go-chat/internal/database"
	"github.com/vladopadikk/go-chat/internal/messages"
	"github.com/vladopadikk/go-chat/internal/user"
	"github.com/vladopadikk/go-chat/internal/ws"
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

	chatRepo := chat.NewRepository(db)
	chatService := chat.NewService(chatRepo)
	chatHandler := chat.NewHandler(chatService)

	messageRepo := messages.NewRepository(db)
	messageService := messages.NewService(messageRepo, chatRepo)
	messageHandler := messages.NewHandler(messageService)

	hub := ws.NewHub()
	go hub.Run()

	wsHandler := ws.NewHandler(hub, chatService, messageService)

	api := router.Group("/api")
	user.RegisterRoutes(api, userHandler)
	auth.RegisterRoutes(api, authHandler)

	protected := api.Group("")
	protected.Use(auth.AuthMiddleware(cfg.JWTSecret))

	chat.RegisterRoutes(protected, chatHandler)
	messages.RegisterRoutes(protected, messageHandler)

	ws.RegisterRoutes(protected, wsHandler)

	router.Run(":" + cfg.AppPort)

}
