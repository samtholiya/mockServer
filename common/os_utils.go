package common

import (
	"os"
)

//GetEnv returns the environment value if available else fallback string
func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Infof("%v not found %v value used", key, fallback)
		return fallback
	}
	return value
}
