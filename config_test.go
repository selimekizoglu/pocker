package main

import (
	"testing"
)

func TestDefaultConfig_consul(t *testing.T) {
	config := DefaultConfig()
	expected := "localhost:8500"
	if config.Consul != expected {
		t.Errorf("expected consul location to be %s but got %s", expected, config.Consul)
	}
}
