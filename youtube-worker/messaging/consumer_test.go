package messaging

import (
	"github.com/segmentio/kafka-go"
	"reflect"
	"testing"
	"time"
)

func Test_of(t *testing.T) {
	type args struct {
		message kafka.Message
	}
	tests := []struct {
		name string
		args args
		want *KafkaMessage
	}{
		{
			name: "simple",
			args: args{
				message: kafka.Message{
					Topic:         "topic",
					Partition:     1,
					Offset:        2,
					HighWaterMark: 3,
					Key:           []byte("key"),
					Value:         []byte("value"),
					Headers:       nil,
					WriterData:    nil,
					Time:          time.Time{},
				},
			},
			want: &KafkaMessage{
				Topic:         "topic",
				Partition:     1,
				Offset:        2,
				HighWaterMark: 3,
				Key:           []byte("key"),
				Value:         []byte("value"),
				Headers:       nil,
				WriterData:    nil,
				Time:          time.Time{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := of(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("of() = %v, want %v", got, tt.want)
			}
		})
	}
}
