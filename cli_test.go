package main

import (
	"testing"
)

func TestParseFlags_consul(t *testing.T) {
	cli := NewCLI()
	config, err := cli.parseFlags([]string{
		"-consul", "10.2.222.32",
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := "10.2.222.32"
	if config.Consul != expected {
		t.Errorf("expected %q to be %q", config.Consul, expected)
	}
}

func TestParseFlags_noConsul(t *testing.T) {
	cli := NewCLI()
	config, err := cli.parseFlags([]string{})
	if err != nil {
		t.Fatal(err)
	}

	expected := "localhost"
	if config.Consul != expected {
		t.Errorf("expected %q to be %q", config.Consul, expected)
	}
}

func TestParseFlags_service(t *testing.T) {
	cli := NewCLI()
	config, err := cli.parseFlags([]string{
		"-service", "customers-api",
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := "customers-api"
	if config.Service != expected {
		t.Errorf("expected %q to be %q", config.Service, expected)
	}
}

func TestParseFlags_noService(t *testing.T) {
	cli := NewCLI()
	config, err := cli.parseFlags([]string{})
	if err != nil {
		t.Fatal(err)
	}

	expected := ""
	if config.Service != "" {
		t.Errorf("expected %q to be %q", config.Service, expected)
	}
}

func TestParseFlags_endpoint(t *testing.T) {
	cli := NewCLI()
	config, err := cli.parseFlags([]string{
		"-endpoint", "/health",
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := "/health"
	if config.Endpoint != expected {
		t.Errorf("expected %q to be %q", config.Endpoint, expected)
	}
}

func TestParseFlags_noEndpoint(t *testing.T) {
	cli := NewCLI()
	config, err := cli.parseFlags([]string{})
	if err != nil {
		t.Fatal(err)
	}

	expected := "/"
	if config.Endpoint != expected {
		t.Errorf("expected %q to be %q", config.Endpoint, expected)
	}
}

func TestParseFlags_expect(t *testing.T) {
	cli := NewCLI()
	config, err := cli.parseFlags([]string{
		"-expect", "3",
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := 3
	if config.Expect != expected {
		t.Errorf("expected %d to be %d", config.Expect, expected)
	}
}

func TestParseFlags_expectAtLeast(t *testing.T) {
	cli := NewCLI()
	config, err := cli.parseFlags([]string{
		"-expectAtLeast", "3",
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := 4
	if expected >= config.Expect {
		t.Errorf("expected %d must higher than %d", expected, config.Expect)
	}
}

func TestParseFlags_retry(t *testing.T) {
	cli := NewCLI()
	_, err := cli.parseFlags([]string{
		"-retry", "3",
	})
	if err != nil {
		t.Fatal(err)
	}
}
