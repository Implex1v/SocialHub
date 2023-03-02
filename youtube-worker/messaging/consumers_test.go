package messaging

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewYoutubePollConsumer(t *testing.T) {
	got := NewYoutubePollConsumer()

	assert.Equal(t, "youtube-poll-consumer", got.Name)
	assert.Equal(t, "youtube-worker", got.GroupId)
	assert.Equal(t, "youtube.poll", got.Topic)
	assert.Equal(t, "localhost:19092", got.BrokerUrl)
}
