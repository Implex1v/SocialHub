package logging

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

func NewSugarLogger(log *zap.Logger) *zap.SugaredLogger {
	return log.Sugar()
}

var Module = fx.Module("logging",
	fx.Provide(
		NewLogger,
		NewSugarLogger,
	),
)
