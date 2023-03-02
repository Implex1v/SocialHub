package metrics

import (
	"github.com/labstack/echo-contrib/prometheus"
	"go.uber.org/fx"
)

type AppMetrics struct {
	KafkaConsumeSum   *prometheus.Metric
	KafkaConsumeCount *prometheus.Metric
}

func NewMetrics() *AppMetrics {
	return &AppMetrics{
		KafkaConsumeSum: &prometheus.Metric{
			Name:        "kakfa_consume_time",
			Description: "details about consumed messages",
			Type:        "summary_vec",
			ID:          "kakfa_consume_time",
			Args:        []string{"topic", "consumer_group", "success"},
		},
		KafkaConsumeCount: &prometheus.Metric{
			Name:        "kafka_consume_count",
			Description: "number of consumed messages",
			Type:        "counter_vec",
			ID:          "kafka_consume_count",
			Args:        []string{"topic", "consumer_group"},
		},
	}
}

func (m *AppMetrics) ToList() []*prometheus.Metric {
	return []*prometheus.Metric{
		m.KafkaConsumeSum,
		m.KafkaConsumeCount,
	}
}

var Module = fx.Module("metrics",
	fx.Provide(
		NewMetrics,
	),
)
