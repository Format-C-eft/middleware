package api

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/Format-C-eft/middleware/internal/database/cache"
	"github.com/Format-C-eft/middleware/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

// OtherMetods - execution of the rest of any methods
func OtherMetods(ginContext *gin.Context) {

	span, _ := opentracing.StartSpanFromContext(ginContext.Request.Context(), "OtherMetods")
	defer span.Finish()

	SessionInfo := ginContext.MustGet("SessionInfo").(*config.SessionInfo)
	cfg := config.GetConfigInstance()

	defer ginContext.Request.Body.Close()

	bodyByte, _ := ioutil.ReadAll(ginContext.Request.Body)
	bodyReader := bytes.NewReader(bodyByte)

	clientHTTP := *http.DefaultClient
	clientHTTP.Timeout = cfg.Servers.OneC.MaxTimeout

query:
	request, err := NewRequest(SessionInfo, bodyReader, ginContext)
	if err != nil {
		logger.WarnKV(ginContext.Request.Context(), "Сould not create a request", "err", err)
		GenerateErrorResponse("Сould not create a request", "Error internal", http.StatusInternalServerError, ginContext)
		ginContext.Abort()
		return
	}

	response, err := clientHTTP.Do(request)
	if err != nil {
		logger.ErrorKV(ginContext.Request.Context(), "Server not responding", "err", err)
		GenerateErrorResponse("Service is unavailable", "Service is unavailable", http.StatusRequestTimeout, ginContext)
		return
	}

	if response.StatusCode == http.StatusNotFound && SessionInfo.Session.Cookie != "" {
		SessionInfo.Session.Cookie = ""
		span.Finish()
		span, _ = opentracing.StartSpanFromContext(ginContext.Request.Context(), "OtherMetods")
		goto query
	}

	defer response.Body.Close()

	respBodyByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.ErrorKV(ginContext.Request.Context(), "Error getting response body", "err", err)
		GenerateErrorResponse("Error getting response body", "The request failed", response.StatusCode, ginContext)
		return
	}

	if response.StatusCode != http.StatusOK {
		if errCheck := checkResponseStatusCode(response.StatusCode, string(respBodyByte), ginContext); errCheck != nil {
			ginContext.Abort()
			return
		}
	}

	ContentType := response.Header.Get("Content-Type")
	if ContentType != "" {
		ginContext.Header("Content-Type", ContentType)
	}

	ginContext.String(response.StatusCode, string(respBodyByte))

	cookies := response.Cookies() // обновляем куки на всякий случай
	for _, cookie := range cookies {
		if cookie.Name == "ibsession" {
			SessionInfo.Session.Cookie = cookie.Value
		}
	}

	err = cache.ClientDB.RefreshExpire(SessionInfo, true)
	if err != nil {
		logger.WarnKV(ginContext.Request.Context(), "CacheDB save error", "err", err)
		GenerateErrorResponse("Error saving the current session", "Error save session", http.StatusInternalServerError, ginContext)
		ginContext.Abort()
		return
	}
}
