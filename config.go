package main

type Config struct {
	// Consul is the address of the consul instance to query
	Consul string

	// Service is name of the service to be poked
	Service string

	// Endpoint is the endpoint of the service to ve poked
	Endpoint string

	// Expect is the number of registered service instances
	Expect int
}

func DefaultConfig() *Config {
	return &Config{
		Consul: "localhost:8500",
		Expect: 1,
	}
}
