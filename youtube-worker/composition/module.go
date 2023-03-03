package composition

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"youtube-worker/config"
	"youtube-worker/http"
	"youtube-worker/logging"
	"youtube-worker/messaging"
	"youtube-worker/metrics"
)

var RootModule = fx.Options(
	config.Module,
	logging.Module,
	metrics.Module,
	http.Module,
	messaging.Module,
	fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: log}
	}),
)
