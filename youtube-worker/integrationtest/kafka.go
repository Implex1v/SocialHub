package integrationtest

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/fx"
	"os"
	"testing"
	"time"
	"youtube-worker/composition"
	"youtube-worker/config"
)

type IntegrationBlock func()

// see: https://golang.testcontainers.org/quickstart/

func IntegrationTest(t *testing.T, block IntegrationBlock) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "docker.redpanda.com/vectorized/redpanda:v22.3.11",
		ExposedPorts: []string{"9092:9092/tcp"},
		Cmd:          []string{"redpanda", "start"},
		WaitingFor:   wait.ForLog("Initialized cluster_id to"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Error(err)
	}

	mappedPort, err := container.MappedPort(ctx, "9092")
	if err != nil {
		t.Error("failed to get port", err)
	}

	err = os.Setenv("KAFKA_PORT", mappedPort.Port())
	if err != nil {
		t.Error("failed tp set env KAFKA_PORT", err)
	}

	app := fx.New(
		composition.RootModule,
	)

	go func() {
		app.Run()
	}()
	defer app.Done()

	block()

	if err := container.Terminate(ctx); err != nil {
		t.Fatalf("failed to terminate kafka: %v", err.Error())
	}
}

func NewMessage(t *testing.T, topic string, message string) {
	conf := config.NewKafkaConfig()
	conn, err := kafka.DialLeader(context.Background(), "tcp", conf.Host+":"+conf.Port, topic, 0)
	if err != nil {
		t.Fatal("failed to dial leader:", err)
	}

	time.Sleep(2 * time.Second)
	_, err = conn.WriteMessages(kafka.Message{Value: []byte(message)})
	if err != nil {
		t.Fatal("failed to write message", err, topic, message)
	}
}
