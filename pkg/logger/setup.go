package logger

import (
	"go.uber.org/zap"
)

// Setup initializes a zap logger for DI.
func Setup() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return l, nil
}
