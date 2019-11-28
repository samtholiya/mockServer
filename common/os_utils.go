package common

import (
	"fmt"
	"os"
)

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		fmt.Printf("%v not found %v value used", key, fallback)
		return fallback
	}
	return value
}
