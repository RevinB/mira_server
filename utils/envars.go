package utils

import (
	"os"
	"strconv"
)

// GetenvInt panics if error occurs
func GetenvInt(key string) int {
	e := os.Getenv(key)
	ret, err := strconv.Atoi(e)
	if err != nil {
		panic("environment variable is not of int type: " + key)
	}

	return ret
}

// GetenvByteArray
func GetenvByteArray(key string) []byte {
	return []byte(os.Getenv(key))
}
