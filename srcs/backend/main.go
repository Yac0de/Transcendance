package main

import (
	"api/controllers"
	"api/database"
	"api/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	database.New()

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(config))
	router.Use(middleware.Token())
	router.Static("/users/avatar", "./avatars")

	users := router.Group("/users")
	auth := router.Group("/auth")

	users.Use(middleware.AuthGuard())

	controllers.Auth(auth)
	controllers.Users(users)

	router.Run(":4000")
}
