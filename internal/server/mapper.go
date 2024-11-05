package server

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/mishakrpv/kafka-proxy-producer/internal/config"
)

type upstreamRoute struct {
	path    string
	methods []string
	params  []param
	tprt    *kafka.TopicPartition
}

func mapRoutes(cfg *config.ProxyConfig) []upstreamRoute {
	routes := []upstreamRoute{}

	for _, route := range cfg.Routes {
		routes = append(routes, upstreamRoute{
			path:    route.UpstreamPathTemplate,
			methods: route.UpstreamHTTPMethod,
			params:  extractParams(route.DownstreamMessage, nil, nil),
			tprt:    route.DownstreamTopicPartition,
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
