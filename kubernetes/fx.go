package kubernetes

import "go.uber.org/fx"

var Module = fx.Module("kubernetes", fx.Provide(NewKubernetesClient))
