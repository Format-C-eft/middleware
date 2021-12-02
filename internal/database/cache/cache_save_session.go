package cache

import (
	"encoding/json"

	"github.com/Format-C-eft/middleware/internal/config"
)

func (c *clientStuct) SaveSession(session *config.SessionInfo) error {

	nameKeys := getNamesKeys(session.Сredentials.Login, session.Session.ID)

	pipe := c.Client.Pipeline().Pipeline()

	pipe.Set(c.Client.Context(), nameKeys.SessionUserKey, nameKeys.Users, c.Expires.Session)
	pipe.Set(c.Client.Context(), nameKeys.SessionCookie1C, session.Session.Cookie, c.Expires.Cookie1C)

	sessionDesc := session.GetAgentInfo()
	sessionDesc.CreateTime = session.Info.CreateTime
	sessionDesc.LastTime = config.NewCurrentTime()
	sessionDesc.EndsTime.Time = config.NewCurrentTime().Time.Add(c.Expires.Session)

	sessionDescByte, err := json.Marshal(sessionDesc)
	if err != nil {
		return err
	}

	pipe.HSet(c.Client.Context(), nameKeys.UsersSessionList, session.Session.ID, string(sessionDescByte))
	pipe.Expire(c.Client.Context(), nameKeys.UsersSessionList, c.Expires.User)

	pipe.HSet(c.Client.Context(), nameKeys.Users, "Login", session.Сredentials.Login)
	pipe.HSet(c.Client.Context(), nameKeys.Users, "Password", session.Сredentials.Password)
	pipe.Expire(c.Client.Context(), nameKeys.Users, c.Expires.User)

	if _, err = pipe.Exec(c.Client.Context()); err != nil {
		return err
	}

	return nil
}
