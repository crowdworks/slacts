// Package config defines what environment variables should be defined for each objects
package config

import (
	"os"
)

// getStringEnv returns environment variable
// if the key doesn't exist, return default value
func getStringEnv(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}
