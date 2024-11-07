package message

import (
	"testing"
)

func TestMakeMessage(t *testing.T) {
	msg := MakeMessage()
	if len(msg) != 0 {
		t.Errorf("expected empty message, got %v", msg)
	}
}

func TestAdd(t *testing.T) {
	msg := MakeMessage()
	msg.Add("key1", "value1").Add("key2", 123)

	if val, ok := msg["key1"]; !ok || val != "value1" {
		t.Errorf("expected key1 to have value 'value1', got %v", val)
	}

	if val, ok := msg["key2"]; !ok || val != 123 {
		t.Errorf("expected key2 to have value 123, got %v", val)
	}
}

func TestBuild(t *testing.T) {
	msg := MakeMessage()
	msg.Add("key1", 123).Add("key2", MakeMessage().Add("key3", "value1"))

	expected := `{"key1":123,"key2":{"key3":"value1"}}`
	result := msg.Build()

	if expected != result {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}
