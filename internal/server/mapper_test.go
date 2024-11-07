package server

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/mishakrpv/kafka-proxy-producer/internal/config"
	"github.com/mishakrpv/kafka-proxy-producer/message"
	"github.com/stretchr/testify/assert"
)

func TestParamWriteValue(t *testing.T) {
	msg := message.MakeMessage()
	p := &param{
		keys:    []string{"outer", "inner"},
		message: msg,
	}
	err := p.writeValue("test-value")

	assert.NoError(t, err, "expected no error from writeValue")
	assert.Equal(t, "test-value", msg["outer"].(map[string]interface{})["inner"], "expected correct nested value")
}

func TestParamKey(t *testing.T) {
	p := &param{keys: []string{"outer", "inner", "key"}}
	key, err := p.key()

	assert.NoError(t, err, "expected no error from key")
	assert.Equal(t, "key", key, "expected correct last key")
}

func TestMapRoutes(t *testing.T) {
	topic := "example-topic"
	cfg := &config.ProxyConfig{
		Routes: []config.Route{
			{
				UpstreamPathTemplate:     "/example",
				UpstreamHTTPMethod:       []string{"GET"},
				DownstreamMessage:        map[string]interface{}{"outer": map[string]interface{}{"key": "[fromquery]"}},
				DownstreamTopicPartition: &kafka.TopicPartition{Topic: &topic},
			},
		},
	}

	routes := mapRoutes(cfg)

	assert.Len(t, routes, 1, "expected one route")
	assert.Equal(t, "/example", routes[0].path, "expected correct path")
	assert.Contains(t, routes[0].methods, "GET", "expected correct method")
	assert.Equal(t, "example-topic", *routes[0].topicPartition.Topic, "expected correct topic partition")
	assert.Equal(t, "[fromquery]", routes[0].params[0].source, "expected correct source in params")
	assert.Equal(t, "key", routes[0].params[0].keys[1], "expected correct key in params")
}

func TestExtractParams(t *testing.T) {
	json := map[string]interface{}{
		"outer": map[string]interface{}{
			"key1": "[fromquery]",
			"key2": "[fromheader]",
		},
	}

	params := extractParams(json, nil, nil)

	assert.Len(t, params, 2, "expected two params")
	assert.Equal(t, "[fromquery]", params[0].source, "expected source from first param")
	assert.Equal(t, "key1", params[0].keys[1], "expected correct first key")
	assert.Equal(t, "[fromheader]", params[1].source, "expected source from second param")
	assert.Equal(t, "key2", params[1].keys[1], "expected correct second key")
}
