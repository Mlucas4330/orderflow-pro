// internal/middleware/prometheus.go
package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "orderflow_http_request_duration_seconds",
	Help: "Duração de todas as requisições HTTP.",
}, []string{"path", "method", "status_code"})

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		statusCode := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()
		method := c.Request.Method

		httpDuration.WithLabelValues(path, method, statusCode).Observe(duration.Seconds())
	}
}
