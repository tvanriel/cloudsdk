//go:build linux
// +build linux

package downwardsapi

import "go.uber.org/fx"

var Module = fx.Module("downwardsapi", fx.Provide(NewDownwardsAPI))
