package messaging

import (
	"github.com/segmentio/kafka-go"
	"time"
)

type KafkaConsumer struct {
	BrokerUrl string
	Topic     string
	GroupId   string
	Name      string
	consume   func(message *KafkaMessage) error
}

type KafkaMessage struct {
	Topic         string
	Partition     int
	Offset        int64
	HighWaterMark int64
	Key           []byte
	Value         []byte
	Headers       []kafka.Header
	WriterData    interface{}
	Time          time.Time
}

func of(message kafka.Message) *KafkaMessage {
	return &KafkaMessage{
		Topic:         message.Topic,
		Partition:     message.Partition,
		Offset:        message.Offset,
		HighWaterMark: message.HighWaterMark,
		Key:           message.Key,
		Value:         message.Value,
		Headers:       message.Headers,
		WriterData:    message.WriterData,
		Time:          message.Time,
	}
}
