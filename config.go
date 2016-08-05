package main

import (
	"github.com/selimekizoglu/gotry"
)

type Config struct {
	// Consul is the address of the consul instance to query
	Consul string

	// Service is name of the service to be poked
	Service string

	// Endpoint is the endpoint of the service to ve poked
	Endpoint string

	// Expect is the number of registered service instances
	Expect int

	//ExpectAtLeast is the minimum number of registered service instances. Service instances can be more.
	ExpectAtLeast int

	//Retry is the number of retries to poke a service in case of failure
	Retry gotry.Retry
}

func DefaultConfig() *Config {
	return &Config{
		Consul: "localhost:8500",
		Expect: 1,
	}
}
