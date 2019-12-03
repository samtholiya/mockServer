package common

import (
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
}

//GetLogger returns the logger for the entire repo
func GetLogger() *logrus.Logger {
	return log
}
