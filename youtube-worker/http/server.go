package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"youtube-worker/config"
)

func NewHTTPServer(lc fx.Lifecycle, logger *zap.SugaredLogger, conf *config.HttpConfig) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				port := ":" + conf.Port
				err := e.Start(port)
				if err != nil {
					logger.Errorf("failed to start server on %v: '%v'", port, err)
					os.Exit(1)
				} else {
					logger.Infof("started http server on %v", port)
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
	routeModule,
	fx.Provide(
		NewHTTPServer,
	),
	fx.Invoke(
		NewPrometheus,
		NewLogger,
		NewRequestLogger,
	),
)
