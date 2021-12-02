package cache

func (c *clientStuct) DropSession(uid, login string) error {

	nameKeys := getNamesKeys(login, uid)

	KeyExistsSession := c.Client.Exists(c.Client.Context(), nameKeys.SessionUserKey)
	KeyExistsCookie := c.Client.Exists(c.Client.Context(), nameKeys.SessionCookie1C)
	KeyExistsSessionList := c.Client.HExists(c.Client.Context(), nameKeys.UsersSessionList, uid)

	if KeyExistsSession.Val() == 1 {
		if err := c.Client.Del(c.Client.Context(), nameKeys.SessionUserKey).Err(); err != nil {
			return err
		}
	}

	if KeyExistsCookie.Val() == 1 {
		if err := c.Client.Del(c.Client.Context(), nameKeys.SessionCookie1C).Err(); err != nil {
			return err
		}
	}

	if KeyExistsSessionList.Val() {
		pipe := c.Client.Pipeline()
		pipe.HDel(c.Client.Context(), nameKeys.UsersSessionList, uid)
		pipe.Expire(c.Client.Context(), nameKeys.UsersSessionList, c.Expires.User)
		pipe.Expire(c.Client.Context(), nameKeys.Users, c.Expires.User)

		if _, err := pipe.Exec(c.Client.Context()); err != nil {
			return err
		}
	}

	return nil
}
