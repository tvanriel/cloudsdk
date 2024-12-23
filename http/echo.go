package http

import (
	"context"
	"regexp"
	"strings"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo-contrib/echoprometheus"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type Http struct {
	Address string
	Engine  *echo.Echo
}

func (h *Http) Use(middleware echo.MiddlewareFunc) {
	h.Engine.Use(middleware)
}

func NewEcho(config Configuration, log *zap.Logger, lc fx.Lifecycle) *Http {
	server := echo.New()

	server.HideBanner = true

	if config.Debug {
		server.Debug = true
	}

	server.Use(echozap.ZapLogger(log.Named("http")))
	server.Use(echoprometheus.NewMiddlewareWithConfig(
		echoprometheus.MiddlewareConfig{},
	))
	if config.Ratelimit > 0 {
		server.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(config.Ratelimit))))
	}

	http := &Http{
		Address: config.Address,
		Engine:  server,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if config.TLS != nil {
				go func() {
					http.Engine.Logger.Fatal(http.Engine.StartTLS(http.Address, config.TLS.CertFile, config.TLS.KeyFile))
				}()
				return nil
			}
			go func() {
				http.Engine.Logger.Fatal(http.Engine.Start(http.Address))
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return http.Engine.Shutdown(ctx)
		},
	})
	return http
}

func RegisterRoutes(routes []RouteGroup, server *Http) {
	for i := range routes {
		routes[i].Handler(server.Engine.Group(makeApiRoute(routes[i])))
	}
}

var prefixSlashes = regexp.MustCompile(`^/+`)

func makeApiRoute(route RouteGroup) string {
	return prefixSlashes.ReplaceAllString(strings.Join([]string{
		"/",
		route.ApiGroup(),
		"/",
		route.Version(),
	}, ""), "/")
}

func (h *Http) EnableDebugging() {
	h.Engine.Debug = true
}

func (h *Http) DisableDebugging() {
	h.Engine.Debug = false
}
