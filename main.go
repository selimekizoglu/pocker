package main

import (
	"os"
)

func main() {
	cli := NewCLI()
	exitCode := cli.Run(os.Args[1:])
	os.Exit(exitCode)
}
