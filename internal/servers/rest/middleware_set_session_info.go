package rest

import (
	"net"
	"net/http"

	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/Format-C-eft/middleware/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
)

// AddSessionInfo - add session info to gin context
func AddSessionInfo() gin.HandlerFunc {
	return func(ginContext *gin.Context) {

		// This will work if the service is hidden behind Apache or nginx and in the settings
		// they have configured to add a header with the real IP of the request
		sessionInfo := config.SessionInfo{}
		RealIP := ginContext.Request.Header.Get("X-Real-IP")

		if net.ParseIP(RealIP) != nil {
			sessionInfo.Info.IP = RealIP
		} else {
			sessionInfo.Info.IP = ginContext.ClientIP()
		}

		stringUA := ginContext.Request.UserAgent()
		structUA := user_agent.New(stringUA)

		if structUA.Bot() {
			// If we could determine that this is a bot, we knock out the session
			logger.InfoKV(ginContext.Request.Context(), "Find bot", "UserAgent", stringUA)
			ginContext.AbortWithStatus(http.StatusNotFound)
			return
		}

		nameEngine, versionEngine := structUA.Engine()
		sessionInfo.Info.Engine = nameEngine + " version: " + versionEngine
		sessionInfo.Info.Browser, sessionInfo.Info.BrowserVersion = structUA.Browser()
		sessionInfo.Info.Mobile = structUA.Mobile()
		sessionInfo.Info.Platform = structUA.Platform()
		sessionInfo.Info.OS = structUA.OS()
		sessionInfo.Info.CreateTime = config.NewCurrentTime()
		sessionInfo.Info.LastTime = config.NewCurrentTime()

		ginContext.Set("SessionInfo", &sessionInfo)

	}
}
