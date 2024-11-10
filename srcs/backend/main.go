package main

import (
	"api/controllers"
	"api/database"
	"api/middleware"
	"api/prometheus"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
)

func main() {
	router := gin.Default()
	database.New()

	// Configuration CORS
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))
	router.Use(middleware.Token())
	router.Use(prometheus.PrometheusMiddleware())

	// Routes
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.Static("/users/avatar", "./avatars")

	users := router.Group("/users")
	auth := router.Group("/auth")
	conversation := router.Group("/conversation")
	controllers.Conversation(conversation)
	users.Use(middleware.AuthGuard())
	controllers.Auth(auth)
	controllers.Users(users)

	router.Run(":4000")
}

