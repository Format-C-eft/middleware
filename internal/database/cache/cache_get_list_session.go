package cache

import (
	"context"
	"encoding/json"

	"github.com/Format-C-eft/middleware/internal/logger"
)

func (c *clientStuct) GetListSession(login, uid string, activeOnly bool) (*[]SessionDescription, error) {

	nameKeys := getNamesKeys(login, "")
	result := []SessionDescription{}

	allKeys := c.Client.HGetAll(c.Client.Context(), nameKeys.UsersSessionList)

	if allKeys.Err() != nil {
		return &result, allKeys.Err()
	}

	for key, val := range allKeys.Val() {
		if uid != "" && uid != key {
			continue
		}

		newItem := SessionDescription{
			SessionID: key,
		}

		if errUnm := json.Unmarshal([]byte(val), &newItem.Info); errUnm != nil {
			logger.ErrorKV(context.TODO(), "Error Unmarshal", "err", errUnm)
		}

		newItem.CheckActive()

		if activeOnly && !newItem.IsActive {
			continue
		}

		result = append(result, newItem)
	}

	return &result, nil
}
