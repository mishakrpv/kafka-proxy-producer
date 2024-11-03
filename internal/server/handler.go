package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/mux"
	"github.com/mishakrpv/kafka-proxy-producer/internal/message"
	"github.com/mishakrpv/kafka-proxy-producer/kafka/producer"
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

func registerRoutes(routes []upstreamRoute) http.Handler {
	router := mux.NewRouter()

	for _, route := range routes {
		router.HandleFunc(route.path, func(w http.ResponseWriter, r *http.Request) {
			f(w, r, route.methods, route.params, route.tprt)
		})
	}

	return router
}

func f(w http.ResponseWriter, r *http.Request, methods []string, params []param, tprt *kafka.TopicPartition) {
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	producer.Produce(tprt, message.Build())

	w.WriteHeader(http.StatusOK)
}

func matchSource(p *param, r *http.Request) error {
	switch p.source {
	case "[FromBody]":
		return fromBody(p, r)
	case "[FromQuery]":
		return fromQuery(p, r)
	case "[FromRoute]":
		return fromRoute(p, r)
	case "[FromForm]":
		return fromForm(p, r)
	case "[FromHeader]":
		return fromHeader(p, r)
	default:
		log.Printf("Unsupported source: %s", p.source)
	}
	return nil
}
