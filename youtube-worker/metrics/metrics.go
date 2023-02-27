package metrics

import (
	"github.com/labstack/echo-contrib/prometheus"
	"go.uber.org/fx"
)

type AppMetrics struct {
	KafkaConsumeSum *prometheus.Metric
}

func NewMetrics() *AppMetrics {
	return &AppMetrics{
		KafkaConsumeSum: &prometheus.Metric{
			Name:        "kakfa_consumer",
			Description: "details about consumed messages",
			Type:        "summary_vec",
			ID:          "kakfa_consumer",
			Args:        []string{"topic", "consumer_group", "success"},
		},
	}
}

func (m *AppMetrics) ToList() []*prometheus.Metric {
	return []*prometheus.Metric{
		m.KafkaConsumeSum,
	}
}

var Module = fx.Module("metrics",
	fx.Provide(
		NewMetrics,
	),
)
