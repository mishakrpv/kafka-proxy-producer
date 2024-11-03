package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/mishakrpv/kafka-proxy-producer/internal/config"
)

type Server struct {
	server *http.Server

	proxyCfg *config.ProxyConfig

	cfg map[string]interface{}
}

func New(proxyCfg *config.ProxyConfig) *Server {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	s := &Server{proxyCfg: proxyCfg}
	s.cfg = make(map[string]interface{})

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	s.cfg["PORT"] = port

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      registerRoutes(mapRoutes(proxyCfg)),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return s
}

func (s *Server) Run() error {
	log.Printf("Server running at %d", s.cfg["PORT"])
	return s.server.ListenAndServe()
}
