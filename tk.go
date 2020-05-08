package tk

import (
	"os"
)

// Env returns value of environment variable
func Env(key, defaultValue string) string {
	value, defined := os.LookupEnv(key)
	if !defined {
		return defaultValue
	}

	return value
}
