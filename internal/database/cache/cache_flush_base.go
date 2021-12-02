package cache

func (c *clientStuct) FlushBase() error {

	if result := c.Client.FlushDB(c.Client.Context()); result.Err() != nil {
		return result.Err()
	}

	return nil
}
