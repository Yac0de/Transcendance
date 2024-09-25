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
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(config))
	router.Use(middleware.Token())

	users := router.Group("/users")
	auth := router.Group("/auth")

	users.Use(middleware.AuthGuard())

	controllers.Auth(auth)
	controllers.Users(users)

	router.Run(":4000")
}
