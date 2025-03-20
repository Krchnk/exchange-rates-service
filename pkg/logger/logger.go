package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	if lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		logger.SetLevel(lvl)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
	return logger
}
