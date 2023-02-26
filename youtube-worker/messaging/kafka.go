package messaging

import "go.uber.org/fx"

var Module = fx.Module("messaging",
	consumerModule,
)
