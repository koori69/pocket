// Package pocket K·J Create at 2020-04-09 21:36
package pocket

import "github.com/sirupsen/logrus"

var logger *logrus.Logger

func init() {
	logger = logrus.StandardLogger()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.AddHook(NewContextHook())
}

func GetLogger() *logrus.Logger {
	return logger
}
