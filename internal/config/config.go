package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ProxyConfig struct {
	Routes        []route       `json:"Routes"`
	LauchSettings lauchSettings `json:"LaunchSettings"`
}

type route struct {
	DownstreamTopicPartition *kafka.TopicPartition  `json:"DownstreamTopicPartition"`
	DownstreamMessage        map[string]interface{} `json:"DownstreamMessage"`
	UpstreamPathTemplate     string                 `json:"UpstreamPathTemplate"`
	UpstreamHTTPMethod       []string               `json:"UpstreamHttpMethod"`
}

type lauchSettings struct {
	Port                 int               `json:"Port"`
	EnvironmentVariables map[string]string `json:"EnvironmentVariables"`
}

func LoadFromFile(path string) *ProxyConfig {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		log.Fatal("An error occurred while opening the configuration file")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("An error occurred while reading the configuration file")
	}

	var config *ProxyConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("An error occurred while reading the configuration file")
	}

	return config
}
