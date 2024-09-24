package main

import (
	"api/controllers"
	"api/database"

	"github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    database.New()
    controllers.Users(router.Group("/users"))
    router.Run(":4000")
}
