package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/selimekizoglu/gotry"
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

func (p *Pocker) Poke() (int, error) {
	consulConf := api.DefaultConfig()
	consulConf.Address = p.Config.Consul
	client, err := api.NewClient(consulConf)
	if err != nil {
		return ExitCodeFail, err
	}

	service := p.Config.Service
	log.Printf("Retrieving service %s from consul (%s)", service, p.Config.Consul)

	// Fetch services from consul catalog
	services, _, err := client.Catalog().Service(service, "", &api.QueryOptions{})
	if err != nil {
		return ExitCodeConsulError, err
	}

	numServices := len(services)
	if numServices != p.Config.Expect {
		return ExitCodeFail, fmt.Errorf("Expected %d services, but consul returned %d", p.Config.Expect, numServices)
	}

	// Iterate through service instances and try to poke the service
	for _, s := range services {
		url := fmt.Sprintf("http://%s:%d%s", s.ServiceAddress, s.ServicePort, p.Config.Endpoint)

		err := gotry.Try(func() error {
			log.Printf("Poking %s", url)
			resp, err := p.Client.Get(url)
			defer resp.Body.Close()

			if err != nil {
				log.Print(err)
				return err
			}
			if resp.StatusCode > 399 {
				log.Printf("StatusCode: %d", resp.StatusCode)
				return errors.New(fmt.Sprintf("Response has status code: %d", resp.StatusCode))
			}

			return nil
		}, p.Config.Retry)

		if err != nil {
			return ExitCodeFail, err
		}
	}

	return ExitCodeOK, nil
}