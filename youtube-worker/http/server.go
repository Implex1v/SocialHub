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
			err := e.Start(":8000")
			if err != nil {
				fmt.Println("Failed to start server on :8000")
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Close()
		},
	})

	return e
}
