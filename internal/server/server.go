package server

import (
	"net/http"

	"github.com/mishakrpv/kafka-proxy-producer/internal/config"
)

func New(cfg *config.ProxyConfig) *http.Server {
	return &http.Server{}
}
