package logging

import (
	_ "github.com/tvanriel/zaplipgloss"
	"go.uber.org/zap"
)

func NewZapLogger(config Configuration) (*zap.Logger, error) {
	if config.Development {
		cfg := zap.NewDevelopmentConfig()
		cfg.Encoding = "lipgloss"
		return cfg.Build()

	}

	return zap.NewProduction()
}
