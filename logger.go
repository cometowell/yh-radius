package main

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
)

// Lumberjackrus github.com/orandin/lumberjackrus 支持日志切割
// Hook for logging to the local filesystem (with logrotate and a file per log level)
func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)

	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  "logs/info.log",
		logrus.ErrorLevel: "logs/error.log",
	}
	logger.AddHook(lfshook.NewHook(pathMap, &logrus.TextFormatter{}))

	return logger
}