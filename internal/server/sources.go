package server

import (
	"encoding/json"
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
	err = p.writeValue(value)
	if err != nil {
		return nil
	}

	return nil
}

func fromQuery(p *param, r *http.Request) error {
	key, err := p.key()
	if err != nil {
		return err
	}

	value := r.URL.Query().Get(key)
	err = p.writeValue(value)
	if err != nil {
		return err
	}
	return nil
}

func fromRoute(p *param, r *http.Request) error {
	key, err := p.key()
	if err != nil {
		return err
	}

	vars := mux.Vars(r)

	if value, ok := vars[key]; ok {
		err = p.writeValue(value)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("no value provided: %s", key)
}

func fromForm(p *param, r *http.Request) error {
	key, err := p.key()
	if err != nil {
		return err
	}

	r.ParseForm()
	value := r.FormValue(key)
	err = p.writeValue(value)
	if err != nil {
		return err
	}
	return nil
}

func fromHeader(p *param, r *http.Request) error {
	key, err := p.key()
	if err != nil {
		return err
	}

	value := r.Header.Get(key)
	err = p.writeValue(value)
	if err != nil {
		return err
	}
	return nil
}
