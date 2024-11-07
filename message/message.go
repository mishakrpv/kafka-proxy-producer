package message

import (
	"encoding/json"
	"log"
)

type Message map[string]interface{}

func MakeMessage() Message {
	return make(map[string]interface{})
}

func (m *Message) Add(key string, value interface{}) Message {
	(*m)[key] = value
	return *m
}

func (m *Message) Build() string {
	data, err := json.Marshal(m)
	if err != nil {
		log.Println("An error occurred while marshaling JSON:", err)
		return ""
	}

	return string(data)
}
