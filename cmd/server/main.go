package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vladopadikk/go-chat/internal/config"
	"github.com/vladopadikk/go-chat/internal/database"
	"github.com/vladopadikk/go-chat/internal/user"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)
	defer db.Close()

	router := gin.Default()

	repo := user.NewRepository(db)
	service := user.NewService(repo)
	handler := user.NewHandler(service)

	api := router.Group("/api")
	user.RegisterRoutes(api, handler)

	// auth.RegisterRoutes(api, db)
	// chat.RegisterRoutes(api, db)
	// message.RegisterRoutes(api, db)
	// ws.RegisterRoutes(api, db)

	router.Run(":" + cfg.AppPort)

}
