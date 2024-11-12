package messaging

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"

	"github.com/mishakrpv/kafka-proxy-producer/internal/config"
	"github.com/mishakrpv/kafka-proxy-producer/internal/server"
)

func TestWithKafka(t *testing.T) {
	ctx := context.Background()

	kafkaContainer, err := kafka.Run(ctx,
		"confluentinc/confluent-local:7.5.0",
		kafka.WithClusterID("test-cluster"),
		testcontainers.WithHostPortAccess(9093),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(kafkaContainer); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		t.Logf("failed to start container: %s", err)
		return
	}

	brokers, err := kafkaContainer.Brokers(ctx)
	require.NoError(t, err)

	t.Logf("Kafka brokers: %s", strings.Join(brokers, " | "))

	os.Setenv("KAFKA__BOOTSTRAP_SERVERS", brokers[0])

	cfg := config.LoadFromFile("../../configuration.json")

	server := server.New(cfg)
	go server.Run()

	sendRequest(t)
}

func sendRequest(t *testing.T) {
	payload := map[string]interface{}{
		"principal": map[string]string{
			"name": "John Doe",
		},
	}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "http://localhost:5465/items/123?page=3", bytes.NewBuffer(payloadBytes))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer your_token")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
