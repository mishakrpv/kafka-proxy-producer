package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gorilla/mux"
)

func isSourceSupported(source string) bool {
	supportedSources := []string{
		"[frombody]",
		"[fromquery]",
		"[fromroute]",
		"[fromform]",
		"[fromheader]",
	}

	return slices.Contains(supportedSources, strings.ToLower(source))
}

func searchInBody(body interface{}, keys []string) string {
	if len(keys) == 0 {
		return ""
	}

	currentKey := keys[0]

	if m, ok := body.(map[string]interface{}); ok {
		if value, exists := m[currentKey]; exists {
			if len(keys) == 1 {
				return fmt.Sprintf("%v", value)
			}
			return searchInBody(value, keys[1:])
		}
	}

	return ""
}

func fromBody(p *param, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	var body map[string]interface{}

	err := decoder.Decode(&body)
	if err != nil {
		return err
	}

	value := searchInBody(body, p.keys)
	p.writeValue(value)
	return nil
}

func fromQuery(p *param, r *http.Request) error {
	key, err := p.key()
	if err != nil {
		return err
	}

	value := r.URL.Query().Get(key)
	p.writeValue(value)
	return nil
}

func fromRoute(p *param, r *http.Request) error {
	key, err := p.key()
	if err != nil {
		return err
	}

	vars := mux.Vars(r)

	if value, ok := vars[key]; ok {
		p.writeValue(value)
		return nil
	}

	return errors.New("no value provided")
}

func fromForm(p *param, r *http.Request) error {
	key, err := p.key()
	if err != nil {
		return err
	}

	r.ParseForm()
	value := r.FormValue(key)
	p.writeValue(value)
	return nil
}

func fromHeader(p *param, r *http.Request) error {
	key, err := p.key()
	if err != nil {
		return err
	}

	value := r.Header.Get(key)
	p.writeValue(value)
	return nil
}
