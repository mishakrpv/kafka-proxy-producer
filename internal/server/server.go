package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/mishakrpv/kafka-proxy-producer/internal/config"
	"github.com/mishakrpv/kafka-proxy-producer/producer"
)

type Server struct {
	server *http.Server

	proxyCfg *config.ProxyConfig

	producer producer.Producer
}

const (
	defaultPort = 5465
)

func New(proxyCfg *config.ProxyConfig) *Server {
	s := &Server{proxyCfg: proxyCfg}

	if proxyCfg.LauchSettings.Port == 0 {
		proxyCfg.LauchSettings.Port = defaultPort
	}

	for key, value := range proxyCfg.LauchSettings.EnvironmentVariables {
		err := os.Setenv(key, value)
		if err != nil {
			log.Fatalf("An error occurred while setting environment variable for: %s", key)
		}
	}

	s.producer = producer.New(kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA__BOOTSTRAP_SERVERS"),
	})

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", proxyCfg.LauchSettings.Port),
		Handler:      s.registerRoutes(mapRoutes(proxyCfg)),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return s
}

func (s *Server) Run() error {
	log.Printf("Server listening on port: %d", s.proxyCfg.LauchSettings.Port)
	return s.server.ListenAndServe()
}
