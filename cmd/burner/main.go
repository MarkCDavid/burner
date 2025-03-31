package main

import (
	"flag"
	// "path/filepath"

	"github.com/MarkCDavid/burner/internal"
)

const (
	FLAG_CONFIGURATION_PATH = "configuration"
	FLAG_SEED               = "seed"
	FLAG_RUNS               = "runs"
)

func main() {
	// configuration_path := flag.String(FLAG_CONFIGURATION_PATH, "./configuration/default.yaml", "Path to the yaml configuration for the simulator")
	seed := flag.Int64(FLAG_SEED, 0, "Seed for the simulation")
	runs := flag.Int64(FLAG_RUNS, 1, "Run count")
	flag.Parse()

	if !flagProvided(FLAG_SEED) {
		seed = nil
	}

	// absolute_configuration_path, err := filepath.Abs(*configuration_path)
	// if err != nil {
	// 	panic(1)
	// }

	for i := int64(0); i < *runs; i++ {
		// internal.Simulate(absolute_configuration_path, seed)
		internal.PoB(seed)
	}
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
