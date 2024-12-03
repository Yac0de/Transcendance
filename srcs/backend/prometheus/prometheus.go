package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/gin-gonic/gin"
	"time"
)

var (

	httpRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
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

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		
		c.Next()
		duration := time.Since(start)
		status := c.Writer.Status()

	
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