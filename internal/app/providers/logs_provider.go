package providers

import (
	"go.uber.org/zap"
)

func LogsProvider() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}
