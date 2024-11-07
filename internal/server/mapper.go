package server

import (
	"errors"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/mishakrpv/kafka-proxy-producer/internal/config"
	"github.com/mishakrpv/kafka-proxy-producer/message"
)

type param struct {
	source  string
	keys    []string
	message message.Message
}

func (p *param) writeValue(value string) error {
	depth := len(p.keys) - 1
	key := p.keys[depth]

	currentMap := p.message

	for i := 0; i <= depth; i++ {
		if i == depth {
			currentMap.Add(key, value)
			return nil
		}

		currentKey := p.keys[i]

		if _, ok := currentMap[currentKey]; !ok {
			currentMap.Add(currentKey, make(map[string]interface{}))
		}
		if m, ok := currentMap[currentKey].(map[string]interface{}); ok {
			currentMap = m
		} else {
			return fmt.Errorf("no value provided: %s", key)
		}
	}

	return nil
}

func (p *param) key() (string, error) {
	length := len(p.keys)
	if length < 1 {
		return "", errors.New("params has no keys")
	}
	return p.keys[length-1], nil
}

type upstreamRoute struct {
	path           string
	methods        []string
	params         []param
	topicPartition *kafka.TopicPartition
}

func mapRoutes(cfg *config.ProxyConfig) []upstreamRoute {
	routes := []upstreamRoute{}

	for _, route := range cfg.Routes {
		routes = append(routes, upstreamRoute{
			path:           route.UpstreamPathTemplate,
			methods:        route.UpstreamHTTPMethod,
			params:         extractParams(route.DownstreamMessage, nil, nil),
			topicPartition: route.DownstreamTopicPartition,
		})
	}
	return routes
}

func extractParams(json map[string]interface{}, keys []string, params *[]param) []param {
	if keys == nil {
		keys = []string{}
	}
	if params == nil {
		params = &[]param{}
	}

	for key, value := range json {
		switch source := value.(type) {
		case string:
			if !isSourceSupported(source) {
				break
			}
			param := &param{source: source, keys: keys}
			param.keys = append(param.keys, key)
			*params = append(*params, *param)
		case map[string]interface{}:
			ks := make([]string, len(keys))
			copy(ks, keys)
			ks = append(ks, key)
			v := value.(map[string]interface{})

			extractParams(v, ks, params)
		}
	}

	return *params
}
