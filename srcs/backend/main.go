package main

import (
	"api/controllers"
	"api/database"
	"api/middleware"
	"api/prometheus"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	router := gin.Default()
	database.New()

	// Configuration CORS
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:4000", "http://localhost:3000", "http://localhost:5173", "http://localhost:8000"},
		// UNCOMMENT FOR PROD MODE AND COMMENT THE ONE ABOVE
		// AllowOrigins:     []string{"*", "http://localhost:4000", "http://localhost:3000", "http://localhost:5173", "http://localhost:8000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))
	router.Use(func(c *gin.Context) {
		c.Set("db", database.DB)
		c.Next()
	})
	router.Use(middleware.Token())
	router.Use(prometheus.PrometheusMiddleware())

	// Routes
	router.Static("/users/avatar", "./avatars")
	router.POST("/api/game-history", controllers.SaveGameHistory)
	router.GET("/api/game-history/:nickname", controllers.GetUserGameHistory)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	users := router.Group("/users")
	auth := router.Group("/auth")
	conversation := router.Group("/conversation")
	controllers.Conversation(conversation)
	users.Use(middleware.AuthGuard())
	controllers.Auth(auth)
	controllers.Users(users)

	router.Run(":4000")
}
