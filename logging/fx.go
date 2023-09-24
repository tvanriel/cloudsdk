package logging

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var Module = fx.Module("logging", fx.Provide(NewZapLogger))

func FXLogger() fx.Option {
	return fx.WithLogger(func(l *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: l}
	})

}
