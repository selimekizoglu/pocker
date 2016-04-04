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

func TestDefaultConfig_expected(t *testing.T) {
	config := DefaultConfig()
	expected := 1
	if config.Expect != expected {
		t.Errorf("expected expect to be %d but got %d", expected, config.Expect)
	}
}
