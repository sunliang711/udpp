package utils

import (
	"github.com/sirupsen/logrus"
)

var lvlMap = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"error": logrus.ErrorLevel,
	"Info":  logrus.InfoLevel,
	"fatal": logrus.FatalLevel,
	"panic": logrus.PanicLevel,
}

//LogLevel translate string to logrus.Level
func LogLevel(l string) logrus.Level {
	if l, ok := lvlMap[l]; ok {
		return l
	}
	return logrus.InfoLevel
}
