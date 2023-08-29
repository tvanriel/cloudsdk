package logging

import "go.uber.org/zap"

func NewZapLogger(config Configuration) (*zap.Logger, error) {
        if config.Development {
                return zap.NewDevelopment()
        }

        return zap.NewProduction()
}
