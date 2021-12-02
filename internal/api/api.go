package api

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/Format-C-eft/middleware/internal/logger"
	"github.com/gin-gonic/gin"
)

// global errors
var (
	errorAuth    = errors.New("error auth")
	errorConnect = errors.New("error connect")
	errorLocked  = errors.New("startup of infobase session is not allowed")
)

// NewRequest - create new request
func NewRequest(SessionInfo *config.SessionInfo, body io.Reader, ginContext *gin.Context) (request1C *http.Request, err error) {

	cfg := config.GetConfigInstance()

	url := cfg.Servers.OneC.Path + ginContext.Request.URL.RequestURI()

	request1C, err = http.NewRequestWithContext(
		ginContext.Request.Context(),
		ginContext.Request.Method,
		url,
		body,
	)
	if err != nil {
		logger.WarnKV(ginContext.Request.Context(), "Request creation error", "err", err)
		return
	}

	if SessionInfo.Session.Cookie != "" {
		request1C.Header.Add("Cookie", "ibsession="+SessionInfo.Session.Cookie)
	} else if SessionInfo.Сredentials.Login != "" {
		request1C.SetBasicAuth(SessionInfo.Сredentials.Login, SessionInfo.Сredentials.Password)
		request1C.Header.Add("IBSession", "start")
	} else {
		err = errors.New("login password and cookies are empty, request headers are not formed")
	}

	request1C.RemoteAddr = ginContext.Request.RemoteAddr
	request1C.Header.Set("User-Agent", ginContext.Request.UserAgent())

	ContentType := ginContext.GetHeader("Content-Type")
	if ContentType != "" {
		request1C.Header.Set("Content-Type", ContentType)
	}

	return
}

// SendRequest - send request
func SendRequest(SessionInfo *config.SessionInfo, ginContext *gin.Context) (*http.Response, string, error) {

	defer ginContext.Request.Body.Close()

	bodyByte, _ := ioutil.ReadAll(ginContext.Request.Body)
	bodyReader := bytes.NewReader(bodyByte)

	request, err := NewRequest(SessionInfo, bodyReader, ginContext)

	if err != nil {
		return &http.Response{}, "", err
	}

	cfg := config.GetConfigInstance()

	client := http.DefaultClient
	client.Timeout = cfg.Servers.OneC.MaxTimeout
	response, err := client.Do(request)
	if err != nil {
		logger.ErrorKV(ginContext.Request.Context(), "The server is not responding", "err", err)
		return &http.Response{}, "", errorConnect
	}

	// Если 404 и мы выполняли запрос с куками значит кто-то прибил сессию на стороне 1С
	// Значит надо попробовать выполнить с базовой авторизацией
	if response.StatusCode == http.StatusNotFound && SessionInfo.Session.Cookie != "" {

		SessionInfo.Session.Cookie = ""

		request, err = NewRequest(SessionInfo, bodyReader, ginContext)
		if err != nil {
			return &http.Response{}, "", err
		}

		response, err = client.Do(request)
		if err != nil {
			logger.ErrorKV(ginContext.Request.Context(), "the server is not responding", "err", err)
			return &http.Response{}, "", errorConnect
		}
	}

	defer response.Body.Close()

	respBodyByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.ErrorKV(ginContext.Request.Context(), "Error getting response body", "err", err)
		GenerateErrorResponse("Error getting response body", "The request failed", response.StatusCode, ginContext)
		return &http.Response{}, "", errors.New("error getting response body")
	}

	return response, string(respBodyByte), err
}

func checkResponseStatusCode(StatusCode int, body string, ginContext *gin.Context) (err error) {

	if StatusCode == http.StatusUnauthorized {
		// The session is live, but the login or password has changed
		err = errorAuth
	} else if StatusCode == http.StatusForbidden && strings.Contains(body, "Startup of infobase session is not allowed") {
		// Session blocking is set and the session is still alive
		err = errorLocked
	} else if StatusCode == http.StatusForbidden {
		// No rights
		err = errorAuth
	} else if StatusCode == http.StatusNotFound {
		// There is a publication, but the server is off
		err = errorConnect
	} else if StatusCode == http.StatusBadGateway {
		// Server 1c is off
		err = errorConnect
	} else if StatusCode == http.StatusInternalServerError {
		// The session is not live and sessions are locked
		if strings.Contains(body, "Startup of infobase session is not allowed") {
			err = errorLocked
		} else if strings.Contains(body, "Database server not found") {
			err = errorConnect
		}
	} else if StatusCode == http.StatusServiceUnavailable {
		// DBMS server responds, but no DB
		if strings.Contains(body, "Database server not found") {
			err = errorConnect
		}
	}

	switch err {
	case errorAuth:
		GenerateErrorResponse("Authorization error Invalid username or password", "Error not a valid login/pass", http.StatusUnauthorized, ginContext)
	case errorConnect:
		GenerateErrorResponse("Service is unavailable", "Service is unavailable", http.StatusRequestTimeout, ginContext)
	case errorLocked:
		messagelock := strings.ReplaceAll(body, "\xEF\xBB\xBF", "") // Removing a BOM character from a string
		messagelock = strings.ReplaceAll(messagelock, "Startup of infobase session is not allowed.", "")
		if len(messagelock) != 0 {
			messagelock = strings.Trim(messagelock, "\r\n")
			messagelock = strings.Trim(messagelock, "\n")
		} else {
			messagelock = body
		}

		GenerateErrorResponse(messagelock, "Service locked", http.StatusLocked, ginContext)
	}

	if err != nil {
		logger.InfoKV(
			ginContext.Request.Context(),
			"checkResponseStatusCode",
			"err", err,
			"StatusCode", StatusCode,
			"body", body,
		)
	}

	return err
}

func convertResponseToRequest(response *http.Response, ginContext *gin.Context, bodyString string) {

	ContentType := response.Header.Get("Content-Type")
	if ContentType != "" {
		ginContext.Header("Content-Type", ContentType)
	}

	ginContext.String(response.StatusCode, bodyString)
}

func findCookie(response *http.Response) (Cookie string) {

	cookies := response.Cookies()
	for _, cookie := range cookies {
		if strings.Compare(cookie.Name, "ibsession") == 0 {
			Cookie = cookie.Value
			return
		}
	}

	return
}
