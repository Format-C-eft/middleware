package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/Format-C-eft/middleware/internal/database/cache"
	"github.com/Format-C-eft/middleware/internal/logger"
	"github.com/Format-C-eft/middleware/internal/token"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

// Login - method login
func Login(ginContext *gin.Context) {

	SessionInfo := ginContext.MustGet("SessionInfo").(*config.SessionInfo)
	cfg := config.GetConfigInstance()
	ctx := ginContext.Request.Context()

	Login, Password, ok := ginContext.Request.BasicAuth()
	if !ok {
		logger.ErrorKV(ctx, "Authorization header is invalid or not passed")
		GenerateErrorResponse("Authorization header is invalid or not passed", "Error auth", http.StatusUnauthorized, ginContext)
		ginContext.Abort()
		return
	}

	SessionInfo.Сredentials.Login = Login
	SessionInfo.Сredentials.Password = Password

	defer ginContext.Request.Body.Close()

	bodyByte, _ := ioutil.ReadAll(ginContext.Request.Body)
	bodyReader := bytes.NewReader(bodyByte)

	request, err := NewRequest(SessionInfo, bodyReader, ginContext)
	if err != nil {
		logger.WarnKV(ctx, "Сould not create a request", "err", err)
		GenerateErrorResponse("Сould not create a request", "Error internal", http.StatusInternalServerError, ginContext)
		ginContext.Abort()
		return
	}

	clientHTTP := *http.DefaultClient
	clientHTTP.Timeout = cfg.Servers.OneC.MaxTimeout

	span, _ := opentracing.StartSpanFromContext(ctx, "Authorization on the 1C side")

	response, err := clientHTTP.Do(request)
	if err != nil {
		span.Finish()
		logger.ErrorKV(ctx, "Server not responding", "err", err)
		GenerateErrorResponse("Service is unavailable", "Service is unavailable", http.StatusRequestTimeout, ginContext)
		return
	}
	span.Finish()

	defer response.Body.Close()

	span, ctx = opentracing.StartSpanFromContext(ctx, "Generate response")
	defer span.Finish()

	respBodyByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.ErrorKV(ctx, "Error getting response body", "err", err)
		GenerateErrorResponse("Error getting response body", "The request failed", response.StatusCode, ginContext)
		return
	}

	if response.StatusCode != http.StatusOK {
		if errCheck := checkResponseStatusCode(response.StatusCode, string(respBodyByte), ginContext); errCheck != nil {
			ginContext.Abort()
			return
		}
	}

	tockenString, sessionUID := token.CreateToken(cfg.Project.Token.Password)

	bodyMap := make(map[string]interface{})
	if errUnm := json.Unmarshal(respBodyByte, &bodyMap); errUnm != nil {
		logger.WarnKV(ctx, "Wrong response from the server", "err", errUnm)
		GenerateErrorResponse("Wrong response from the server", "Server error", http.StatusBadRequest, ginContext)
		return
	}

	if _, ok = bodyMap["data"]; ok {
		bodyMap["data"].(map[string]interface{})["token"] = tockenString
	} else {
		logger.WarnKV(ctx, "Wrong response from the server", "err", err)
		GenerateErrorResponse("Wrong response from the server", "Server error", http.StatusBadRequest, ginContext)
		return
	}

	SessionInfo.Session.ID = sessionUID

	err = cache.ClientDB.SaveSession(SessionInfo)
	if err != nil {
		logger.ErrorKV(ctx, "Cache save error", "err", err)
		GenerateErrorResponse("Session data saving error", "Error save session", http.StatusInternalServerError, ginContext)
		ginContext.Abort()
		return
	}

	ginContext.JSON(http.StatusOK, bodyMap)
}

// LogOut - close session
func LogOut(ginContext *gin.Context) {

	SessionInfo := ginContext.MustGet("SessionInfo").(*config.SessionInfo)

	everyWhere := false

	if strings.Compare(strings.ToLower(ginContext.Query("everywhere")), "true") == 0 {
		everyWhere = true
	}

	ctx := ginContext.Request.Context()

	if everyWhere {
		sessionList, err := cache.ClientDB.GetListSession(SessionInfo.Сredentials.Login, SessionInfo.Session.ID, true)
		if err != nil {
			logger.WarnKV(ctx, "Error GetListSession", "err", err)
			GenerateErrorResponse("Session deletion error", "Error clean session", http.StatusInternalServerError, ginContext)
			ginContext.Abort()
			return
		}

		for _, session := range *sessionList {
			err := cache.ClientDB.DropSession(session.SessionID, SessionInfo.Сredentials.Login)
			if err != nil {
				logger.WarnKV(ctx, "cache clean error", "err", err)
				GenerateErrorResponse("Session deletion error", "Error clean session", http.StatusInternalServerError, ginContext)
				ginContext.Abort()
				return
			}
		}
	} else {
		err := cache.ClientDB.DropSession(SessionInfo.Session.ID, SessionInfo.Сredentials.Login)
		if err != nil {
			logger.ErrorKV(ctx, "cache clean error", "err", err)
			GenerateErrorResponse("Session deletion error", "Error clean session", http.StatusInternalServerError, ginContext)
			ginContext.Abort()
			return
		}
	}
	GenerateErrorResponse("", "SUCCESS", http.StatusOK, ginContext)
}
