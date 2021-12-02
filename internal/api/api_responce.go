package api

import (
	"github.com/gin-gonic/gin"
)

type descriptionErr struct {
	Data struct {
	} `json:"data"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// GenerateErrorResponse - error response generation procedure
func GenerateErrorResponse(textErr, codeErr string, errStatusCode int, ginContext *gin.Context) {

	ginContext.JSON(errStatusCode,
		descriptionErr{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{codeErr, textErr},
		},
	)

}
