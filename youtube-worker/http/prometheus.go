package http

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"youtube-worker/metrics"
)

func NewPrometheus(echo *echo.Echo, metrics *metrics.AppMetrics) {
	p := prometheus.NewPrometheus("echo", nil, metrics.ToList())
	p.Use(echo)
}
