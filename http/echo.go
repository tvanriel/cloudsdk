package http

import (
	"strings"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Http struct {
        Address string
        engine *echo.Echo
}

func (h *Http) Use(middleware echo.MiddlewareFunc) {
        h.engine.Use(middleware)
}

func NewEcho(config Configuration, log *zap.Logger) (*Http) {
        server := echo.New()

        server.Use(echozap.ZapLogger(log))
        server.Use(echoprometheus.NewMiddlewareWithConfig(
                echoprometheus.MiddlewareConfig{},
        ))

        
        
        return &Http{
                Address: config.Address,
                engine: server,
        }
}

func RegisterRoutes(routes []RouteGroup, server *Http) {
        for i := range routes {
                routes[i].Handler(server.engine.Group(makeApiRoute(routes[i])))
        }
}

func makeApiRoute(route RouteGroup) string {
        return strings.Join([]string{
                "/",
                route.ApiGroup(),
                "/",
                route.Version(),
        }, "")
}



func Listen(server *Http) {
        go server.engine.Start(server.Address)
}
