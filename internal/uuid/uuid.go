package uuid

import (
	"github.com/google/uuid"
	"strings"
)

func GetUuidFromPathSecondPosition(path string) uuid.UUID {
	return getUuidFromPathIndex(path, 2)
}

func GetUuidFromPathThirdPosition(path string) uuid.UUID {
	return getUuidFromPathIndex(path, 3)
}

func getUuidFromPathIndex(path string, index int) uuid.UUID {
	parts := strings.Split(path, "/")
	return uuid.MustParse(parts[index])
}
