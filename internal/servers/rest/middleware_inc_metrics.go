package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	methods = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "middleware",
		Help: "method counter",
	}, []string{"name", "method"})
)

// IncMetrics - inc metrics
func IncMetrics() gin.HandlerFunc {
	return func(ginContext *gin.Context) {

		methods.WithLabelValues(ginContext.Request.URL.Path, ginContext.Request.Method).Inc()

	}
}
