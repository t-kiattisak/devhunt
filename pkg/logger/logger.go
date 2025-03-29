package logger

import (
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Init() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
}

func Info(message string) {
	log.Info(message)
}

func Error(message string) {
	log.Error(message)
}

func Fatal(message string) {
	log.Fatal(message)
}
