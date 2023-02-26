package http

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"net/http"
)

// see: https://echo.labstack.com/guide/

type status struct {
	Status string
}

func Status(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, status{Status: "Ok"})
	})
}

var RouteModule = fx.Module(
	"routes",
	fx.Invoke(
		Status,
	),
)
