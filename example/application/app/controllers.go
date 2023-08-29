package app

import (
        "github.com/tvanriel/cloudsdk/http"
        "github.com/labstack/echo/v4"
)


type TestController struct {}

func (t *TestController) Version() string {
        return "v1"
}

func (t *TestController) Handler(g *echo.Group) {
        g.GET("/", func(c echo.Context) error {
                c.String(200, "dawai dawai!")
                return nil
        })
}

func (t *TestController) ApiGroup() string {
        return "test"
}

func NewHttpController() http.RouteGroup {
        return &TestController{}
}
