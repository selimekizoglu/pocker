package main

import (
	"flag"
	"github.com/selimekizoglu/gotry"
	"log"
	"time"
)

type PockerConfigCallback func(*Config)

// CLI is the main entry point for Pocker
type CLI struct {
}

func NewCLI() (cli *CLI) {
	return &CLI{}
}

func (cli *CLI) Run(args []string) int {
	log.Printf("Running: %s\n", args)
	conf, err := cli.parseFlags(args)
	if err != nil {
		return ExitCodeParseFlagsError
	}

	pocker := NewPocker(conf)
	exitCode, _ := pocker.Poke()

	return exitCode
}

// parseFlags is a helper function for parsing command line flags
func (cli *CLI) parseFlags(args []string) (*Config, error) {
	flags := flag.NewFlagSet("pocker", flag.ExitOnError)

	consul := flags.String("consul", "localhost", "")
	service := flags.String("service", "", "")
	endpoint := flags.String("endpoint", "/", "")
	expect := flags.Int("expect", 1, "")
	retry := flags.Int("retry", 0, "")
	timeout := flags.Int("timeout", 0, "")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	config := &Config{
		Consul:   *consul,
		Service:  *service,
		Endpoint: *endpoint,
		Expect:   *expect,
		Retry: &gotry.Retry{
			Max:     uint(*retry),
			Timeout: time.Duration(*timeout) * time.Millisecond,
		},
	}

	return config, nil
}
