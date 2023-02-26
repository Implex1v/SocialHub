package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"youtube-worker/http"
	"youtube-worker/logging"
	"youtube-worker/messaging"
)

func main() {
	fx.New(
		module,
	).Run()
}

var module = fx.Options(
	http.Module,
	logging.Module,
	messaging.Module,
	fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: log}
	}),
)
