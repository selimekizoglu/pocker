package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/selimekizoglu/gotry"
	"log"
	"net/http"
)

const (
	ExitCodeOK int = 0

	ExitCodeParseFlagsError = 10 + iota
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

	var services []*api.CatalogService
	err = gotry.Try(func() error {
		// Fetch services from consul catalog
		services, _, err = client.Catalog().Service(service, "", &api.QueryOptions{})
		if err != nil {
			log.Print(err)
			return err
		}
		numServices := len(services)
		if numServices != p.Config.Expect {
			err = fmt.Errorf("Expected %d services, but consul returned %d", p.Config.Expect, numServices)
			log.Print(err)
			return err
		}

		return nil
	}, p.Config.Retry)

	if err != nil {
		return ExitCodeFail, err
	}

	// Iterate through service instances and try to poke the service
	for _, s := range services {
		url := fmt.Sprintf("http://%s:%d%s", s.ServiceAddress, s.ServicePort, p.Config.Endpoint)

		err := gotry.Try(func() error {
			log.Printf("Poking %s", url)
			resp, err := p.Client.Get(url)
			if err != nil {
				log.Print(err)
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode > 399 {
				err = fmt.Errorf("Response has status code: %d", resp.StatusCode)
				log.Print(err)
				return err
			}

			return nil
		}, p.Config.Retry)

		if err != nil {
			return ExitCodeFail, err
		}
	}

	return ExitCodeOK, nil
}
