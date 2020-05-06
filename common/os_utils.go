package common

import (
	"crypto/rand"
	"fmt"
	"os"
)

//GetEnv returns the environment value if available else fallback string
func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Debugf("%v not found %v value used", key, fallback)
		return fallback
	}
	return value
}

//GetUniqueString returns string of length
func GetUniqueString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	return s
}

//CreateFolder creates folder if it does not exists
func CreateFolder(folderName string) {
	_, err := os.Stat(folderName)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folderName, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
}
