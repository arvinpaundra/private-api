package util

import "github.com/google/uuid"

func GenerateUUID() string {
	return uuid.Must(uuid.NewV7()).String()
}

func ParseUUID(id string) uuid.UUID {
	return uuid.MustParse(id)
}
