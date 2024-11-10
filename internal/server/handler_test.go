package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/mishakrpv/kafka-proxy-producer/message"
	"github.com/stretchr/testify/assert"
)

type mockProducer struct {
}

func (m *mockProducer) Produce(topicPartition *kafka.TopicPartition, msg string) error {
	return nil
}

func newMockServer() *Server {
	return &Server{
		producer: &mockProducer{},
	}
}

func TestRegisterRoutes(t *testing.T) {
	server := newMockServer()
	routes := []upstreamRoute{
		{path: "/test", methods: []string{"GET"}},
	}

	router := server.registerRoutes(routes)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "expected status OK for registered route")
}

func TestMatchSource(t *testing.T) {
	tests := []struct {
		source     string
		setupReq   func() *http.Request
		expected   string
		expectFail bool
	}{
		{
			source: "[FromQuery]",
			setupReq: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/?key=queryValue", nil)
			},
			expected: "queryValue",
		},
		{
			source: "[FromHeader]",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("key", "headerValue")
				return req
			},
			expected: "headerValue",
		},
		{
			source: "[FromBody]",
			setupReq: func() *http.Request {
				body := `{"key": "bodyValue"}`
				return httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
			},
			expected: "bodyValue",
		},
		{
			source: "[FromForm]",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				req.Form = make(map[string][]string)
				req.Form.Set("key", "formValue")
				return req
			},
			expected: "formValue",
		},
	}

	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			msg := message.MakeMessage()
			p := param{source: test.source, keys: []string{"key"}, message: msg}

			err := matchSource(&p, test.setupReq())

			if test.expectFail {
				assert.Error(t, err, "expected error for invalid source")
			} else {
				assert.NoError(t, err, "expected no error for valid source")
				assert.Equal(t, test.expected, msg["key"], "expected correct value for source")
			}
		})
	}
}
