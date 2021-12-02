package api

import (
	"crypto/sha1" //nolint
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/Format-C-eft/middleware/internal/database/cache"
	"github.com/Format-C-eft/middleware/internal/logger"
	"github.com/gin-gonic/gin"
)

// CheckLogin - check login
func CheckLogin(ginContext *gin.Context) {

	Login, Password, ok := ginContext.Request.BasicAuth()
	if !ok {
		logger.ErrorKV(ginContext.Request.Context(), "Authorization header is invalid or not passed")
		GenerateErrorResponse("Authorization header is invalid or not passed", "Error auth", http.StatusUnauthorized, ginContext)
		ginContext.Abort()
		return
	}

	cfg := config.GetConfigInstance()

	hashFunc := sha1.New() //nolint
	hashFunc.Write([]byte(strings.ToUpper(Password)))
	passBase64 := base64.StdEncoding.EncodeToString(hashFunc.Sum(nil))

	query := ginContext.Request.URL.Query()
	query.Set("login", Login)
	query.Set("passwordHash", passBase64)

	ginContext.Request.URL.RawQuery = query.Encode()
	ginContext.Request.URL.Path = "/check-login"

	sessionInfo, err := cache.ClientDB.GetSessionInfo("check-login")
	if err != nil && err != cache.ErrorNotFound {
		logger.ErrorKV(ginContext.Request.Context(), "Error get session", "err", err)
		GenerateErrorResponse("Internal server error", "Internal server error", http.StatusInternalServerError, ginContext)
		ginContext.Abort()
		return
	}

	sessionInfo.Сredentials.Login = strings.ToLower(cfg.Servers.OneC.User.Login)
	sessionInfo.Сredentials.Password = cfg.Servers.OneC.User.Password
	sessionInfo.Session.ID = "check-login"

	response, bodyString, err := SendRequest(sessionInfo, ginContext)
	if err != nil {
		switch err {
		case errorConnect:
			GenerateErrorResponse("Service is unavailable", "Service is unavailable", http.StatusServiceUnavailable, ginContext)
		default:
			GenerateErrorResponse("Internal server error", "Internal server error", http.StatusInternalServerError, ginContext)
		}

		ginContext.Abort()
		return
	}

	if response.StatusCode != http.StatusOK {
		if errCheck := checkResponseStatusCode(response.StatusCode, bodyString, ginContext); errCheck != nil {
			ginContext.Abort()
			return
		}
	}

	convertResponseToRequest(response, ginContext, bodyString)

	sessionInfo.Session.Cookie = findCookie(response)

	err = cache.ClientDB.RefreshExpire(sessionInfo, true)
	if err != nil {
		logger.WarnKV(ginContext.Request.Context(), "CacheDB save error", "err", err)
		GenerateErrorResponse("Error saving the current session", "Error save session", http.StatusInternalServerError, ginContext)
		ginContext.Abort()
		return
	}
}
