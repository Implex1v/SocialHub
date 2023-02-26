package main

import (
	"go.uber.org/fx"
	"youtube-worker/http"
	"youtube-worker/logging"
)

func main() {
	fx.New(
		module,
	).Run()
}

var module = fx.Options(
	http.Module,
	logging.Module,
)
