package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MetodOrPatchNotFound - stub for methods not described
func MetodOrPatchNotFound(ginContext *gin.Context) {
	GenerateErrorResponse("Metod or patch not found",
		"Error path",
		http.StatusNotFound,
		ginContext)
}
