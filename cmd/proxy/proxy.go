package main

import (
	"flag"
	"log"

	"github.com/mishakrpv/kafka-proxy-producer/internal/config"
	"github.com/mishakrpv/kafka-proxy-producer/internal/server"
)

const (
	DEFAULT_PATH = "../../configuration.json"
)

func main() {
	pathPtr := flag.String("c", DEFAULT_PATH, "Path to the configuration file")
	cfg := config.LoadFromFile(*pathPtr)

	server := server.New(cfg)
	log.Fatal(server.Run())
}
