package main

import (
	"log"
	"os"

	"github.com/mishakrpv/kafka-proxy-producer/internal/config"
	"github.com/mishakrpv/kafka-proxy-producer/internal/server"
)

const (
	defaultPath = "configuration.json"
)

func main() {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = defaultPath
	}

	cfg := config.LoadFromFile(path)

	server := server.New(cfg)
	log.Fatal(server.Run())
}
