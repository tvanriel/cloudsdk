package http

import "go.uber.org/fx"

const GROUP_ROUTES = `group:"httproutes"`

var Module = fx.Module(
	"http",
	fx.Invoke(
		fx.Annotate(
			RegisterRoutes,
			fx.ParamTags(GROUP_ROUTES),
		),
	),
	fx.Provide(
		NewEcho,
	),
)

func AsRouteGroup(handler any) any {
	return fx.Annotate(
		handler,
		fx.As(new(RouteGroup)),
		fx.ResultTags(GROUP_ROUTES),
	)
}
