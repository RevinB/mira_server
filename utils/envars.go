package utils

import (
	"os"
)

// GetenvByteArray
func GetenvByteArray(key string) []byte {
	return []byte(os.Getenv(key))
}
