package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	// Counter for the total number of HTTP requests
	httpRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total", // Metric name
			Help: "Total number of HTTP requests", // Metric description
		},
		[]string{"method", "endpoint", "status"}, // Labels for method, endpoint, and status
	)

	// Histogram for measuring the duration of HTTP requests
	httpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds", // Metric name
			Help: "Duration of HTTP requests", // Metric description
		},
		[]string{"method", "endpoint"}, // Labels for method and endpoint
	)

	activeUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Number of currently connected users",
		},
	)

	loginAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "login_attempts_total",
			Help: "Total number of login attempts",
		},
		[]string{"status"},
	)

	gamePlayed = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "total_game_played",
			Help: "Number of games played",
		},
	)
)

// PrometheusMiddleware records metrics for HTTP requests
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Start time of the request
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		
		c.Next() // Call the next handler
		duration := time.Since(start) // Calculate request duration
		status := c.Writer.Status() // Get response status

		// Record metrics
		httpRequests.WithLabelValues(
            c.Request.Method, 
            path, 
            string(rune(status)),
        ).Inc()
        
        httpDuration.WithLabelValues(
            c.Request.Method, 
            path,
        ).Observe(duration.Seconds())
	}
}

func IncrementActiveUsers() {
	activeUsers.Inc()
}

func DecrementActiveUsers() {
	activeUsers.Dec()
}

func RecordLoginAttempt(success bool) {
	if success {
		loginAttempts.WithLabelValues("success").Inc()
	}else{
		loginAttempts.WithLabelValues("failure").Inc()
	}
}

 func IncrementPlayedGames() {
 	gamePlayed.Inc()
}