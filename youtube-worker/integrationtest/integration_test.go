package integrationtest

import (
	"context"
	"fmt"
	"github.com/phayes/freeport"
	"github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/fx"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
	"youtube-worker/composition"
	"youtube-worker/config"
)

type Wrapper struct {
	ctx       context.Context
	container testcontainers.Container
	app       *fx.App
	kafkaPort string
	httpPort  string
}

func (w *Wrapper) RunContainer() error {
	w.ctx = context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "docker.redpanda.com/vectorized/redpanda:v22.3.11",
		ExposedPorts: []string{"9092:9092/tcp"},
		Cmd:          []string{"redpanda", "start"},
		WaitingFor:   wait.ForLog("Initialized cluster_id to"),
	}

	container, err := testcontainers.GenericContainer(w.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	w.container = container
	if err != nil {
		return fmt.Errorf("failed to create container: '%w'", err)
	}

	mappedPort, err := w.container.MappedPort(w.ctx, "9092")
	if err != nil {
		return fmt.Errorf("failed to get kafka port: '%w'", err)
	}

	err = os.Setenv("KAFKA_PORT", mappedPort.Port())
	if err != nil {
		return fmt.Errorf("failed tp set env KAFKA_PORT: '%w'", err)
	}

	httpPort, err := freeport.GetFreePort()
	err = os.Setenv("HTTP_PORT", strconv.Itoa(httpPort))
	if err != nil {
		return fmt.Errorf("failed tp set env HTTP_PORT: '%w'", err)
	}

	w.httpPort = strconv.Itoa(httpPort)
	w.kafkaPort = mappedPort.Port()
	w.app = fx.New(
		composition.RootModule,
	)

	go func() {
		w.app.Run()
	}()

	return nil
}

func (w *Wrapper) CleanUp() {
	if err := w.container.Terminate(w.ctx); err != nil {
		log.Fatalf("failed to terminate kafka: %v", err.Error())
	}
	w.app.Done()
}

func (w *Wrapper) NewMessage(t *testing.T, topic string, message string) {
	conf := config.NewKafkaConfig()
	conn, err := kafka.DialLeader(w.ctx, "tcp", conf.Host+":"+w.kafkaPort, topic, 0)
	if err != nil {
		t.Fatal("failed to dial leader:", err)
	}

	time.Sleep(2 * time.Second)
	_, err = conn.WriteMessages(kafka.Message{Value: []byte(message)})
	if err != nil {
		t.Fatal("failed to write message", err, topic, message)
	}
}
