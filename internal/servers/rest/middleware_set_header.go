package rest

import (
	"time"

	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/gin-gonic/gin"
)

// SetHeaders - set headers for CORS
func SetHeaders() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		cfg := config.GetConfigInstance()

		ginContext.Header("Server", "Apache/2.4.48 (Win64)") // pretending to be an Apache
		ginContext.Header("Date", time.Now().Format(config.LayoutDateFormat))

		ginContext.Header("X-Robots-Tag", "noindex, nofollow") // Directions for search engines not to index

		origin := ginContext.Request.Header.Get("Origin")
		if origin != "" {
			// If the request comes from an allowed source, add headers for CORS to work
			for _, value := range cfg.Services.Rest.AccessOrigin {
				if value == origin {
					ginContext.Header("Access-Control-Allow-Origin", value)
					ginContext.Header("Access-Control-Allow-Credentials", "true")
					ginContext.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Origin, Authorization")
					ginContext.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTION")
					ginContext.Header("Access-Control-Max-Age", "600")
				}
			}
		}
	}
}
