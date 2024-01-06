package prometheus

import "go.uber.org/fx"

var Module = fx.Module("prometheus", 
        fx.Provide(NewPrometheus),
)
