package drivers

import (
	"bytes"
	"encoding/json"
)

func sameLockID(current []byte, requested []byte) bool {
	currentID := lockID(current)
	requestedID := lockID(requested)
	if currentID == "" || requestedID == "" {
		return bytes.Equal(bytes.TrimSpace(current), bytes.TrimSpace(requested))
	}

	return currentID == requestedID
}

func lockID(payload []byte) string {
	var lock struct {
		ID string `json:"ID"`
	}
	if err := json.Unmarshal(payload, &lock); err != nil {
		return ""
	}

	return lock.ID
}
