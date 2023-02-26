package http

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
)

func NewPrometheus(echo *echo.Echo) {
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(echo)
}
