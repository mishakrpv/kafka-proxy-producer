package server

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/mux"
	"github.com/mishakrpv/kafka-proxy-producer/message"
)

func (s *Server) registerRoutes(routes []upstreamRoute) http.Handler {
	router := mux.NewRouter()

	for _, route := range routes {
		router.HandleFunc(route.path, func(w http.ResponseWriter, r *http.Request) {
			s.handleRoute(w, r, route.methods, route.params, route.topicPartition)
		})
	}

	return router
}

func (s *Server) handleRoute(w http.ResponseWriter, r *http.Request, methods []string, params []param, topicPartition *kafka.TopicPartition) {
	for i, method := range methods {
		methods[i] = strings.ToUpper(method)
	}
	if !slices.Contains(methods, r.Method) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	message := message.MakeMessage()

	for _, p := range params {
		p.message = message
		err := matchSource(&p, r)
		if err != nil {
			log.Println("An error occurred matching source:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	err := s.producer.Produce(topicPartition, message.Build())
	if err != nil {
		log.Println("An error occurred producing message:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func matchSource(p *param, r *http.Request) error {
	switch strings.ToLower(p.source) {
	case "[frombody]":
		return fromBody(p, r)
	case "[fromquery]":
		return fromQuery(p, r)
	case "[fromroute]":
		return fromRoute(p, r)
	case "[fromform]":
		return fromForm(p, r)
	case "[fromheader]":
		return fromHeader(p, r)
	default:
		return fmt.Errorf("unsupported source: %s", p.source)
	}
}
