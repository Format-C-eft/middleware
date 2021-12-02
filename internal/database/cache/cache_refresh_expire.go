package cache

import (
	"encoding/json"

	"github.com/Format-C-eft/middleware/internal/config"
)

func (c *clientStuct) RefreshExpire(session *config.SessionInfo, refreshCookie bool) error {

	nameKeys := getNamesKeys(session.Сredentials.Login, session.Session.ID)

	if c.Client.Exists(c.Client.Context(), nameKeys.SessionUserKey).Val() != 1 ||
		c.Client.Exists(c.Client.Context(), nameKeys.Users).Val() != 1 ||
		c.Client.Exists(c.Client.Context(), nameKeys.UsersSessionList).Val() != 1 {
		err := c.SaveSession(session)
		return err
	}

	pipe := c.Client.Pipeline()
	pipe.Expire(c.Client.Context(), nameKeys.SessionUserKey, c.Expires.Session) // Updating the session key lifetime
	pipe.Expire(c.Client.Context(), nameKeys.Users, c.Expires.User)             // Updating the lifetime of the user key

	agentInfo := session.GetAgentInfo()
	agentInfo.EndsTime.Time = agentInfo.LastTime.Time.Add(c.Expires.Session)

	sessionDescByte, err := json.Marshal(agentInfo)
	if err != nil {
		return err
	}
	pipe.HSet(c.Client.Context(), nameKeys.UsersSessionList, session.Session.ID, string(sessionDescByte)) // Update session information
	pipe.Expire(c.Client.Context(), nameKeys.UsersSessionList, c.Expires.User)

	if refreshCookie {
		// We update cookies only if asked, since requests can be not only for 1C
		pipe.Set(c.Client.Context(), nameKeys.SessionCookie1C, session.Session.Cookie, c.Expires.Cookie1C)
	}

	if _, err = pipe.Exec(c.Client.Context()); err != nil {
		return err
	}

	return nil
}
