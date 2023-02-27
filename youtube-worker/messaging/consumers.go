package messaging

import (
	"context"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
	"youtube-worker/metrics"
)

// see: https://github.com/segmentio/kafka-go#reader-

func NewConsumer(l *zap.SugaredLogger, lc fx.Lifecycle, metrics *metrics.AppMetrics) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:19092"},
		Topic:   "youtube.poll",
		GroupID: "youtube-worker",
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				for {
					m, err := r.ReadMessage(context.Background())
					before := time.Now().Local()

					if err != nil {
						labels := prom.Labels{"topic": "youtube.poll", "consumer_group": "youtube-worker", "success": "false"}
						metrics.
							KafkaConsumeSum.
							MetricCollector.(*prom.SummaryVec).With(labels).Observe(time.Now().Local().Sub(before).Seconds())
						l.Errorln("failed to consume message", err)
						break
					}

					l.Infof("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
					labels := prom.Labels{"topic": "youtube.poll", "consumer_group": "youtube-worker", "success": "true"}
					metrics.
						KafkaConsumeSum.
						MetricCollector.(*prom.SummaryVec).With(labels).Observe(time.Now().Local().Sub(before).Seconds())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return r.Close()
		},
	})

}

var consumerModule = fx.Module("consumers", fx.Invoke(NewConsumer))
