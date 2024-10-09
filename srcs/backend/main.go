package main

import (
	"api/controllers"
	"api/database"
	"api/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			help: "Total number of HTTP request",
		},
		[]string{"method", "endpoint", "status"},
	)
	httpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)
)

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Default) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		httpRequests.WithLabelValues(c.Request.Method, c.FullPath(), string(status)).Inc()
		httpDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration.Seconds())
	}
}

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
	router.Use(prometheusMiddleware())
	
	router.GET("/metrics", gin.Wraph(promhttp.Handler()))
	router.Static("/users/avatar", "./avatars")

	users := router.Group("/users")
	auth := router.Group("/auth")

	users.Use(middleware.AuthGuard())

	controllers.Auth(auth)
	controllers.Users(users)

	router.Run(":4000")
}
