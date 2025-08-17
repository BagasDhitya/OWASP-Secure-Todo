package log

import "go.uber.org/zap"

func New() (*zap.Logger, func()) {
	logger, _ := zap.NewProduction() // JSON structured logs
	cleanup := func() { _ = logger.Sync() }
	return logger, cleanup
}
