package uuid

import (
	"github.com/google/uuid"
)

// NewUUID generate uuid string
func New() string {
	return uuid.New().String()
}

// VerifyUUID verify the uuid string is valid
func VerifyUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
