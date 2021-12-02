package cache

import (
	"encoding/json"
	"errors"

	"github.com/Format-C-eft/middleware/internal/config"
)

func (c *clientStuct) GetSessionInfo(uid string) (*config.SessionInfo, error) {

	valUserLink, err := c.Client.Get(c.Client.Context(), "Session:"+uid+":UserKey").Result()
	if err != nil {
		return &config.SessionInfo{}, ErrorNotFound
	}

	SessionInfo := config.SessionInfo{}

	SessionInfo.Session.ID = uid

	valUserInfo, err := c.Client.HGetAll(c.Client.Context(), string(valUserLink)).Result()
	if err != nil {
		return &config.SessionInfo{}, err
	}

	for key, value := range valUserInfo {
		if key == "Login" {
			SessionInfo.Сredentials.Login = value
		} else if key == "Password" {
			SessionInfo.Сredentials.Password = value
		}
	}

	if SessionInfo.Сredentials.Login == "" {
		return &config.SessionInfo{}, errors.New("for key - " + string(valUserLink) + " in base empty login")
	}

	nameKeys := getNamesKeys(SessionInfo.Сredentials.Login, SessionInfo.Session.ID)

	if valSessionCookie, errGet := c.Client.Get(c.Client.Context(), nameKeys.SessionCookie1C).Result(); errGet == nil {
		SessionInfo.Session.Cookie = string(valSessionCookie)
	}

	valAgentInfo, err := c.Client.HGet(c.Client.Context(), nameKeys.UsersSessionList, SessionInfo.Session.ID).Result()
	if err != nil {
		return &config.SessionInfo{}, err
	}

	err = json.Unmarshal([]byte(valAgentInfo), &SessionInfo.Info)
	if err != nil {
		return &config.SessionInfo{}, err
	}

	return &SessionInfo, nil
}
