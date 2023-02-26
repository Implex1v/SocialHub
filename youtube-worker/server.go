package main

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"youtube-worker/http"
)

func main() {
	fx.New(
		fx.Provide(http.NewHTTPServer),
		fx.Invoke(func(echo *echo.Echo) {}),
	).Run()
}
