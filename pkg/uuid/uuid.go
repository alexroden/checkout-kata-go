package uuid

import (
	"strings"

	"github.com/google/uuid"
)

const NilUUID = "00000000-0000-0000-0000-000000000000"

// New creates a new uuid, or returns a nil one for testing.
func New(isTesting bool) string {
	result := NilUUID
	if !isTesting {
		result = uuid.New().String()
	}

	return strings.ToUpper(strings.ReplaceAll(result, "-", ""))
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
