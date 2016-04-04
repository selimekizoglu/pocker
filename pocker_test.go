package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/testutil"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestPoke_healthyService(t *testing.T) {
	consul := testConsul(t)
	defer consul.Stop()

	conf := DefaultConfig()
	conf.Consul = consul.HTTPAddr
	conf.Service = "healthy-service"
	conf.Endpoint = "/health"

	setupConsul(t, consul)
	client := testHTTPClient(t, conf)

	pocker := NewPocker(conf)
	pocker.Client = client
	status := pocker.Poke()
	if status != ExitCodeOK {
		t.Errorf("expected OK but got %d", status)
	}
}

func TestPoke_unhealthyService(t *testing.T) {
	consul := testConsul(t)
	defer consul.Stop()

	conf := DefaultConfig()
	conf.Consul = consul.HTTPAddr
	conf.Service = "unhealthy-service"
	conf.Endpoint = "/health"

	setupConsul(t, consul)
	client := testHTTPClient(t, conf)

	pocker := NewPocker(conf)
	pocker.Client = client
	status := pocker.Poke()
	if status != ExitCodeFail {
		t.Errorf("expected Fail but got %d", status)
	}
}

func TestPoke_noSuchService(t *testing.T) {
	consul := testConsul(t)
	defer consul.Stop()

	conf := DefaultConfig()
	conf.Consul = consul.HTTPAddr
	conf.Service = "unknown-service"
	conf.Endpoint = "/health"

	setupConsul(t, consul)
	client := testHTTPClient(t, conf)

	pocker := NewPocker(conf)
	pocker.Client = client
	status := pocker.Poke()
	if status != ExitCodeFail {
		t.Errorf("expected ConsulFail but got %d", status)
	}
}

func TestPoke_emptyService(t *testing.T) {
	consul := testConsul(t)
	defer consul.Stop()

	conf := DefaultConfig()
	conf.Consul = consul.HTTPAddr
	conf.Service = ""
	conf.Endpoint = "/health"

	setupConsul(t, consul)
	client := testHTTPClient(t, conf)

	pocker := NewPocker(conf)
	pocker.Client = client
	status := pocker.Poke()
	if status != ExitCodeConsulError {
		t.Errorf("expected ConsulFail but got %d", status)
	}
}

func TestPoke_badExpect(t *testing.T) {
	consul := testConsul(t)
	defer consul.Stop()

	conf := DefaultConfig()
	conf.Consul = consul.HTTPAddr
	conf.Service = "healthy-service"
	conf.Endpoint = "/health"
	conf.Expect = 2

	setupConsul(t, consul)
	client := testHTTPClient(t, conf)

	pocker := NewPocker(conf)
	pocker.Client = client
	status := pocker.Poke()
	if status != ExitCodeFail {
		t.Errorf("expected Fail but got %d", status)
	}
}

func TestRun_consulError(t *testing.T) {
	conf := DefaultConfig()
	conf.Service = "healthy-service"
	conf.Endpoint = "/health"

	client := testHTTPClient(t, conf)

	pocker := NewPocker(conf)
	pocker.Client = client
	status := pocker.Poke()
	if status != ExitCodeConsulError {
		t.Errorf("expected ConsulError exit code but got %d", status)
	}
}

type FakeService struct {
	StatusCodes map[string]int
}

func (s *FakeService) RoundTrip(req *http.Request) (*http.Response, error) {
	url := req.URL
	log.Printf("Fakeservice handling %s", req.URL)
	statusCode, ok := s.StatusCodes[url.String()]
	if !ok {
		statusCode = http.StatusNotFound
	}

	body := ioutil.NopCloser(strings.NewReader(""))
	resp := &http.Response{
		StatusCode: statusCode,
		Body:       body,
	}

	return resp, nil
}

func testConsul(t *testing.T) *testutil.TestServer {
	consul := testutil.NewTestServerConfig(t, func(c *testutil.TestServerConfig) {
		c.Stdout = ioutil.Discard
		c.Stderr = ioutil.Discard
	})

	return consul
}

func setupConsul(t *testing.T, server *testutil.TestServer) {
	consulConf := api.DefaultConfig()
	consulConf.Address = server.HTTPAddr
	consul, err := api.NewClient(consulConf)
	if err != nil {
		t.Fatal(err)
	}

	agent := consul.Agent()
	agent.ServiceRegister(&api.AgentServiceRegistration{
		ID:      "healthy-service-1",
		Name:    "healthy-service",
		Port:    8081,
		Address: "localhost",
	})
	agent.ServiceRegister(&api.AgentServiceRegistration{
		ID:      "unhealthy-service-1",
		Name:    "unhealthy-service",
		Port:    8082,
		Address: "localhost",
	})
	agent.ServiceRegister(&api.AgentServiceRegistration{
		ID:      "unhealthy-service-2",
		Name:    "unhealthy-service",
		Port:    8083,
		Address: "localhost",
	})

}

func testHTTPClient(t *testing.T, conf *Config) *http.Client {
	return &http.Client{
		Transport: &FakeService{
			StatusCodes: map[string]int{
				fmt.Sprintf("http://localhost:%d%s", 8081, conf.Endpoint): http.StatusOK,
				fmt.Sprintf("http://localhost:%d%s", 8082, conf.Endpoint): http.StatusOK,
				fmt.Sprintf("http://localhost:%d%s", 8083, conf.Endpoint): http.StatusServiceUnavailable,
			},
		},
	}
}
