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

func NewYoutubePollConsumer() KafkaConsumer {
	return KafkaConsumer{
		BrokerUrl: "localhost:19092",
		Topic:     "youtube.poll",
		GroupId:   "youtube-worker",
		Name:      "youtube-poll-consumer",
		consume: func(message *KafkaMessage) {
			println(string(message.Value))
		},
	}
}

func RegisterConsumers(consumers []KafkaConsumer, l *zap.SugaredLogger, lc fx.Lifecycle, metrics *metrics.AppMetrics) {
	for _, consumer := range consumers {
		l.Debugf("Registering consumer '%v' for topic '%v'", consumer.Name, consumer.Topic)

		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{consumer.BrokerUrl},
			Topic:   consumer.Topic,
			GroupID: consumer.GroupId,
		})

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					for {
						m, err := r.ReadMessage(context.Background())
						go func() {
							before := time.Now().Local()

							if err != nil {
								labels := prom.Labels{
									"topic":          consumer.Topic,
									"consumer_group": consumer.GroupId,
									"success":        "false",
								}
								metrics.
									KafkaConsumeSum.
									MetricCollector.(*prom.SummaryVec).With(labels).Observe(time.Now().Local().Sub(before).Seconds())
								l.Errorln("failed to consume message", err)
								return
							}

							consumer.consume(of(m))
							labels := prom.Labels{
								"topic":          consumer.Topic,
								"consumer_group": consumer.GroupId,
								"success":        "true",
							}
							metrics.
								KafkaConsumeSum.
								MetricCollector.(*prom.SummaryVec).With(labels).Observe(time.Now().Local().Sub(before).Seconds())
						}()
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return r.Close()
			},
		})
	}
}

var consumerModule = fx.Module(
	"consumers",
	fx.Provide(
		fx.Annotate(NewYoutubePollConsumer, fx.ResultTags(`group:"consumers"`)),
	),
	fx.Invoke(
		fx.Annotate(RegisterConsumers, fx.ParamTags(`group:"consumers"`)),
	),
)
