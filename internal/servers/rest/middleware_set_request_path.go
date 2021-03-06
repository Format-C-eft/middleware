package rest

import (
	"strings"

	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/gin-gonic/gin"
)

// ClearPath - clear request path
func ClearPath() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		cfg := config.GetConfigInstance()

		ginContext.Request.URL.Path = strings.Replace(ginContext.Request.URL.Path, cfg.Services.Rest.Path, "", 1)

	}
}
