package messaging

import (
	"context"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"strconv"
	"time"
	"youtube-worker/config"
	"youtube-worker/metrics"
)

// see: https://github.com/segmentio/kafka-go#reader-

func NewYoutubePollConsumer() KafkaConsumer {
	return KafkaConsumer{
		BrokerUrl: "localhost:19092",
		Topic:     "youtube.poll",
		GroupId:   "youtube-worker",
		Name:      "youtube-poll-consumer",
		consume: func(message *KafkaMessage) error {
			println(string(message.Value))
			return nil
		},
	}
}

func RegisterConsumers(consumers []KafkaConsumer, l *zap.SugaredLogger, lc fx.Lifecycle, metrics *metrics.AppMetrics, conf config.KafkaConfig) {
	for _, consumer := range consumers {
		l.Infof("Registering consumer '%v' in consumer group '%v' for topic '%v' on '%v:%v'", consumer.Name, consumer.GroupId, consumer.Topic, conf.Host, conf.Port)

		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{conf.Host + ":" + conf.Port},
			Topic:   consumer.Topic,
			GroupID: consumer.GroupId,
		})

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					for {
						m, err := r.ReadMessage(context.Background())
						if err != nil {
							l.Errorf("failed to consume message (topic='%v', consumer='%v', consumerGroup='%v': '%v'", consumer.Topic, consumer.Name, consumer.GroupId, err)
							break
						}

						metrics.
							KafkaConsumeCount.
							MetricCollector.(*prom.CounterVec).
							With(prom.Labels{"topic": consumer.Topic, "consumer_group": consumer.GroupId}).
							Inc()

						go consume(metrics, consumer, l, m)
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

func consume(appMetrics *metrics.AppMetrics, consumer KafkaConsumer, logger *zap.SugaredLogger, m kafka.Message) {
	before := time.Now().Local()

	err := consumer.consume(of(m))
	if err != nil {
		logger.Errorf("failed to handle message (topic='%v', key='%v', offset='%v', partition='%v', message='%v'): '%v'", m.Topic, string(m.Key), m.Offset, m.Partition, strFirstCharacters(string(m.Value), 500), err)
	} else {
		logger.Debugf("handled message (topic='%v', key='%v', offset='%v', partition='%v', message='%v')", m.Topic, string(m.Key), m.Offset, m.Partition, strFirstCharacters(string(m.Value), 500))
	}

	appMetrics.
		KafkaConsumeSum.
		MetricCollector.(*prom.SummaryVec).
		With(prom.Labels{"topic": consumer.Topic, "consumer_group": consumer.GroupId, "success": strconv.FormatBool(err == nil)}).
		Observe(time.Now().Local().Sub(before).Seconds())
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
