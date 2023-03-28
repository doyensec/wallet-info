package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func InitializeLogger() {
	logger, _ := zap.NewProduction()
	Logger = logger.Sugar()
}
