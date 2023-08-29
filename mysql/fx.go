package mysql

import "go.uber.org/fx"

var Module = fx.Module(
        "mysql",
        fx.Provide(NewGorm),
)
