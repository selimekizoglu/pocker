package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"net/http"
)

const (
	ExitCodeOK int = 0

	ExitCodeParseFlagsError = 10 + iota
	ExitCodeConsulError
	ExitCodeFail
)

type Pocker struct {
	Config *Config

	// HTTPClient to poke services
	Client *http.Client
}

func NewPocker(conf *Config) *Pocker {
	return &Pocker{
		Config: conf,
		Client: &http.Client{},
	}
}

func (p *Pocker) Poke() int {
	consulConf := api.DefaultConfig()
	consulConf.Address = p.Config.Consul
	client, err := api.NewClient(consulConf)
	if err != nil {
		return ExitCodeFail
	}

	service := p.Config.Service
	log.Printf("Retrieving service %s from consul (%s)", service, p.Config.Consul)
	services, _, err := client.Catalog().Service(service, "", &api.QueryOptions{})
	if err != nil {
		return ExitCodeConsulError
	}

	numServices := len(services)
	if numServices != p.Config.Expect {
		return ExitCodeFail
	}

	for _, s := range services {
		url := fmt.Sprintf("http://%s:%d%s", s.ServiceAddress, s.ServicePort, p.Config.Endpoint)
		log.Printf("Poking %s", url)
		resp, err := p.Client.Get(url)
		if err != nil || resp.StatusCode > 399 {
			return ExitCodeFail
		}
		defer resp.Body.Close()
	}

	return ExitCodeOK
}
