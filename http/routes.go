package http

import (
	echo "github.com/labstack/echo/v4"
)

type RouteGroup interface {
        ApiGroup() string
        Version() string
        Handler(*echo.Group)
}


