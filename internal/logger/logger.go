package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func New() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return logger, nil
}
