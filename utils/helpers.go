package utils

import (
	"github.com/google/uuid"
)

func GetTransactionId() string {
	return uuid.Must(uuid.NewRandom()).String()
}
