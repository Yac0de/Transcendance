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
)

// PrometheusMiddleware records metrics for HTTP requests
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Start time of the request
		c.Next() // Call the next handler
		duration := time.Since(start) // Calculate request duration
		status := c.Writer.Status() // Get response status

		// Record metrics
		httpRequests.WithLabelValues(c.Request.Method, c.FullPath(), string(rune(status))).Inc() // Increment request counter
		httpDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration.Seconds()) // Observe request duration
	}
}
