package prometheus

import "go.uber.org/fx"


var Module = fx.Module("prometheus", fx.Invoke(Listen), fx.Provide(NewPrometheus))
