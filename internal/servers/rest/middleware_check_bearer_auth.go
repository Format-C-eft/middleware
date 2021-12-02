package rest

import (
	"net/http"
	"strings"

	"github.com/Format-C-eft/middleware/internal/api"
	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/Format-C-eft/middleware/internal/database/cache"
	"github.com/Format-C-eft/middleware/internal/logger"
	"github.com/Format-C-eft/middleware/internal/token"
	"github.com/gin-gonic/gin"
)

// CheckBearerAuth - check bearer auth
func CheckBearerAuth() gin.HandlerFunc {
	return func(ginContext *gin.Context) {

		cfg := config.GetConfigInstance()

		tokenHeader := ginContext.Request.Header.Get("Authorization")

		if tokenHeader == "" {
			logger.ErrorKV(ginContext.Request.Context(), "Authorization data missing")
			api.GenerateErrorResponse("Authorization data missing", "Error auth", http.StatusUnauthorized, ginContext)
			ginContext.Abort()
			return
		}

		if !strings.Contains(tokenHeader, "Bearer") {
			logger.ErrorKV(ginContext.Request.Context(), "Only Bearer authorization allowed")
			api.GenerateErrorResponse("Only Bearer authorization allowed", "Error auth", http.StatusUnauthorized, ginContext)
			ginContext.Abort()
			return
		}

		splitted := strings.Split(tokenHeader, " ") // Divide the string into two the second part and check
		if len(splitted) != 2 {
			logger.ErrorKV(ginContext.Request.Context(), "Invalid authorization string")
			api.GenerateErrorResponse("Invalid authorization string", "Error auth", http.StatusUnauthorized, ginContext)
			ginContext.Abort()
			return
		}

		tokenPart := splitted[1] // We get the second part of the token

		jwtToken, err := token.CheckValid(tokenPart, cfg.Project.Token.Password)
		if err != nil {
			logger.ErrorKV(ginContext.Request.Context(), "Invalid JWT token", "err", err)
			api.GenerateErrorResponse("Invalid JWT token", "Error auth", http.StatusUnauthorized, ginContext)
			ginContext.Abort()
			return
		}

		// TODO - проверить
		cacheSession, err := cache.ClientDB.GetSessionInfo(jwtToken.UID)
		if err != nil {
			logger.ErrorKV(ginContext.Request.Context(), "Session not found", "err", err)
			api.GenerateErrorResponse("Session not found", "Error auth", http.StatusUnauthorized, ginContext)
			ginContext.Abort()
		}

		currentSession := ginContext.MustGet("SessionInfo").(*config.SessionInfo)

		if cacheSession.Info.Browser != currentSession.Info.Browser ||
			cacheSession.Info.Mobile != currentSession.Info.Mobile ||
			cacheSession.Info.Platform != currentSession.Info.Platform ||
			cacheSession.Info.OS != currentSession.Info.OS {

			logger.ErrorKV(ginContext.Request.Context(), "Сlient changed, session drop")
			api.GenerateErrorResponse("Сlient changed, session drop", "Error auth", http.StatusUnauthorized, ginContext)
			if err := cache.ClientDB.DropSession(cacheSession.Session.ID, cacheSession.Сredentials.Login); err != nil {
				logger.WarnKV(ginContext.Request.Context(), "Auto drop session error", "err", err)
			}

			ginContext.Abort()
			return
		}

		ginContext.Set("SessionInfo", cacheSession)
	}
}
