package main

import (
	"go.uber.org/fx"
	"youtube-worker/http"
)

func main() {
	fx.New(
		fx.Provide(http.NewHTTPServer),
		fx.Invoke(http.NewPrometheus),
	).Run()
}
