package api

import (
	"net/http"
	"strings"

	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/Format-C-eft/middleware/internal/database/cache"
	"github.com/Format-C-eft/middleware/internal/logger"
	"github.com/gin-gonic/gin"
)

type getSessionResponse struct {
	Data struct {
		Items        []cache.SessionDescription `json:"items"`
		ItemsPerPage int                        `json:"itemsPerPage"`
	} `json:"data"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// GetSession - get session
func GetSession(ginContext *gin.Context) {

	ParamUUID := ginContext.Param("UUID")

	activeQuery := ginContext.Query("active")
	activeOnly := false

	if activeQuery != "" && strings.ToLower(activeQuery) == "true" {
		activeOnly = true
	}

	if ParamUUID != "" && !UUIDIsValid(ParamUUID) {
		GenerateErrorResponse("UUID validate error", "ERROR_GET", http.StatusBadRequest, ginContext)
		ginContext.Abort()
		return
	}

	sessionInfo := ginContext.MustGet("SessionInfo").(*config.SessionInfo)

	sessionList, err := cache.ClientDB.GetListSession(sessionInfo.Сredentials.Login, ParamUUID, activeOnly)
	if err != nil {
		logger.WarnKV(ginContext.Request.Context(), "CacheDB error get list session", "err", err)
		GenerateErrorResponse("Error get list session", "Internal server error", http.StatusInternalServerError, ginContext)
		ginContext.Abort()
		return
	}

	response := getSessionResponse{
		Data: struct {
			Items        []cache.SessionDescription "json:\"items\""
			ItemsPerPage int                        "json:\"itemsPerPage\""
		}{
			Items:        *sessionList,
			ItemsPerPage: len(*sessionList),
		},
		Error: struct {
			Code    string "json:\"code\""
			Message string "json:\"message\""
		}{Code: "SUCCESS"},
	}

	ginContext.JSON(http.StatusOK, response)

	err = cache.ClientDB.RefreshExpire(sessionInfo, false)
	if err != nil {
		logger.WarnKV(ginContext.Request.Context(), "CacheDB refresh key error", "err", err)
		GenerateErrorResponse("Error saving the current session", "Error save session", http.StatusInternalServerError, ginContext)
		ginContext.Abort()
		return
	}

}

// DropSession - drop session
func DropSession(ginContext *gin.Context) {

	ParamUUID := ginContext.Param("UUID")

	if ParamUUID != "" && !UUIDIsValid(ParamUUID) {
		GenerateErrorResponse("UUID validate error", "ERROR_GET", http.StatusBadRequest, ginContext)
		ginContext.Abort()
		return
	}

	SessionInfo := ginContext.MustGet("SessionInfo").(*config.SessionInfo)

	if SessionInfo.Session.ID != ParamUUID {

		if err := cache.ClientDB.DropSession(ParamUUID, SessionInfo.Сredentials.Login); err != nil {
			logger.WarnKV(ginContext.Request.Context(), "CacheDB clean error", "err", err)
			GenerateErrorResponse("Error deleting session data", "Error clean session", http.StatusInternalServerError, ginContext)
			ginContext.Abort()
			return
		}

		if err := cache.ClientDB.RefreshExpire(SessionInfo, false); err != nil {
			logger.WarnKV(ginContext.Request.Context(), "CacheDB refresh error", "err", err)
			GenerateErrorResponse("Error saving the current session", "Error refresh session", http.StatusInternalServerError, ginContext)
			ginContext.Abort()
			return
		}
	} else {
		logger.InfoKV(ginContext.Request.Context(), "CacheDB drop current session")
		GenerateErrorResponse("You can not delete the current session", "Error drop session", http.StatusBadRequest, ginContext)
		ginContext.Abort()
		return
	}

	GenerateErrorResponse("", "SUCCESS", http.StatusOK, ginContext)
}
