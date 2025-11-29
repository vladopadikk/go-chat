package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vladopadikk/go-chat/internal/config"
	"github.com/vladopadikk/go-chat/internal/database"
)

/*
func MyLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		println("Request to:", ctx.FullPath())
		ctx.Next()
	}
}
*/

/*
	router.Use(MyLogger())

	api := router.Group("/api")
	{
		api.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"status": "ok"})
		})
		api.GET("/hello", func(ctx *gin.Context) {
			ctx.String(200, "Hello")
		})
		api.GET("/users/:id", func(ctx *gin.Context) {
			id := ctx.Param("id")
			ctx.JSON(200, gin.H{"user_id": id})
		})
		api.GET("/search", func(ctx *gin.Context) {
			query := ctx.Query("q")
			ctx.JSON(200, gin.H{"query": query})
		})
		api.POST("/login", func(ctx *gin.Context) {
			var body struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}

			if err := ctx.BindJSON(&body); err != nil {
				ctx.JSON(400, gin.H{"error": "invalid json"})
				return
			}

			ctx.JSON(200, gin.H{
				"email": body.Email,
			})
		})
	}
*/

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)
	defer db.Close()

	router := gin.Default()

	//api := router.Group("/api")

	// auth.RegisterRoutes(api, db)
	// user.RegisterRoutes(api, db)
	// chat.RegisterRoutes(api, db)
	// message.RegisterRoutes(api, db)
	// ws.RegisterRoutes(api, db)

	router.Run(":" + cfg.AppPort)

}
