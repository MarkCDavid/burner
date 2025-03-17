package main

import (
	"flag"

	"github.com/MarkCDavid/burner/internal"
	"path/filepath"
)

const (
	FLAG_CONFIGURATION_PATH = "configuration"
	FLAG_SEED               = "seed"
)

func main() {
	configuration_path := flag.String(FLAG_CONFIGURATION_PATH, "./configuration/default.yaml", "Path to the yaml configuration for the simulator")
	seed := flag.Int64(FLAG_SEED, 0, "Seed for the simulation")
	flag.Parse()

	if !flagProvided(FLAG_SEED) {
		seed = nil
	}

	absolute_configuration_path, err := filepath.Abs(*configuration_path)
	if err != nil {
		panic(1)
	}

	internal.Simulate(absolute_configuration_path, seed)
}

func flagProvided(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
