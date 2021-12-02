package rest

import (
	"github.com/Format-C-eft/middleware/internal/logger"
	"github.com/gin-gonic/gin"
)

func changeLoggerLevel() gin.HandlerFunc {
	return func(ginContext *gin.Context) {

		loglvl := ginContext.GetHeader("log-level")

		zaplvl, ok := logger.StrToZapLevel(loglvl)
		if !ok {
			return
		}

		ctx := ginContext.Request.Context()

		logger.InfoKV(ctx, "log level change", "levels", zaplvl)
		if parsedLevel, ok := logger.StrToZapLevel(loglvl); ok {
			newLogger := logger.CloneWithLevel(ctx, parsedLevel)
			ctx := logger.AttachLogger(ctx, newLogger)
			ginContext.Request = ginContext.Request.Clone(ctx)
		}

	}
}
