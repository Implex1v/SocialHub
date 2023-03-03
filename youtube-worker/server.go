package main

import (
	"go.uber.org/fx"
	"youtube-worker/composition"
)

func main() {
	fx.New(
		composition.RootModule,
	).Run()
}
