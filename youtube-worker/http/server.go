package http

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func NewHTTPServer(lc fx.Lifecycle) *echo.Echo {
	e := echo.New()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := e.Start(":8000")
				if err != nil {
					fmt.Print(fmt.Errorf("failed to start server on :8000", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Close()
		},
	})

	return e
}

var Module = fx.Module("http",
	RouteModule,
	fx.Provide(
		NewHTTPServer,
	),
	fx.Invoke(
		NewPrometheus,
		NewLogger,
		NewRequestLogger,
	),
)
