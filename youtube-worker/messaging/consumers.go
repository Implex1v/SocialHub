package messaging

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// see: https://github.com/segmentio/kafka-go#reader-

func NewConsumer(l *zap.SugaredLogger, lc fx.Lifecycle) {
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
					if err != nil {
						l.Errorln("failed to consume message", err)
						break
					}

					l.Infof("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
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
