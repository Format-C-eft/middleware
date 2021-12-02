package rest

import (
	"net"
	"time"

	"github.com/Format-C-eft/middleware/internal/logger"
	"github.com/gin-gonic/gin"
)

func ginLogger() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		start := time.Now()

		path := ginContext.Request.URL.Path
		query := ginContext.Request.URL.RawQuery

		xRealIP := ginContext.Request.Header.Get("X-Real-IP")

		realIP := ginContext.ClientIP()

		if net.ParseIP(xRealIP) != nil {
			realIP = xRealIP
		}

		ginContext.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(ginContext.Errors) > 0 {
			// you can add fields if required
			for _, err := range ginContext.Errors.Errors() {
				logger.ErrorKV(ginContext, "errors of context", "err", err)
			}
		} else {
			logger.InfoKV(
				ginContext.Request.Context(),
				"Query info",
				"path", path,
				"query", query,
				"method", ginContext.Request.Method,
				"status", ginContext.Writer.Status(),
				"ip", realIP,
				"latency", latency,
			)
		}
	}
}
