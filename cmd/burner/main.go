package main

import (
	"flag"

	"github.com/MarkCDavid/burner/internal"
)

const (
	FLAG_CONFIGURATION_PATH = "configuration"
)

func main() {
	configuration_path := flag.String(FLAG_CONFIGURATION_PATH, "./configuration/default.yaml", "Path to the yaml configuration for the simulator")

	flag.Parse()

	internal.Simulate(*configuration_path)
}
