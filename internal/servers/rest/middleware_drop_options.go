package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AbortMetodOption - abort method options
func AbortMetodOption() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		if ginContext.Request.Method == http.MethodOptions {
			ginContext.AbortWithStatus(http.StatusOK)
			return
		}
	}
}
