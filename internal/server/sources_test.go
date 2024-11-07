package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mishakrpv/kafka-proxy-producer/message"
	"github.com/stretchr/testify/assert"
)

func newParam(source string, keys []string) *param {
	return &param{
		source:  source,
		keys:    keys,
		message: message.MakeMessage(),
	}
}

func TestFromBody(t *testing.T) {
	bodyJSON := `{"outer":{"inner":"value"}}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	p := newParam("[frombody]", []string{"outer", "inner"})
	err := fromBody(p, req)
	assert.NoError(t, err, "expected no error from fromBody")
	assert.Equal(t, bodyJSON, p.message.Build(), "expected message field to contain the correct value")
}

func TestFromQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?key=value", nil)
	p := newParam("[fromquery]", []string{"key"})

	err := fromQuery(p, req)
	assert.NoError(t, err, "expected no error from fromQuery")
	assert.Equal(t, "value", p.message["key"], "expected message field to contain the correct value")
}

func TestFromRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/items/123", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "123"})

	p := newParam("[fromroute]", []string{"id"})
	err := fromRoute(p, req)
	assert.NoError(t, err, "expected no error from fromRoute")
	assert.Equal(t, "123", p.message["id"], "expected message field to contain the correct route parameter value")
}

func TestFromForm(t *testing.T) {
	formData := "key=value"
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	p := newParam("[fromform]", []string{"key"})
	err := fromForm(p, req)
	assert.NoError(t, err, "expected no error from fromForm")
	assert.Equal(t, "value", p.message["key"], "expected message field to contain the correct form parameter value")
}

func TestFromHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Custom-Header", "header-value")

	p := newParam("[fromheader]", []string{"X-Custom-Header"})
	err := fromHeader(p, req)
	assert.NoError(t, err, "expected no error from fromHeader")
	assert.Equal(t, "header-value", p.message["X-Custom-Header"], "expected message field to contain the correct header value")
}
