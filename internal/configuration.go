package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Seed *int64 `yaml:"seed"`

	NodeCount      int64   `yaml:"node_count"`
	SimulationTime float64 `yaml:"simulation_time_in_seconds"`

	ChainReogranizationThreshold   int64   `yaml:"chain_reorganization_threshold"`
	AverageNetworkLatencyInSeconds float64 `yaml:"average_network_latency_in_seconds"`

	ProofOfWork ProofOfBurnConfiguration `yaml:"proof_of_work"`
}

type ProofOfBurnConfiguration struct {
	Enabled                        bool    `yaml:"enabled"`
	EpochLength                    int64   `yaml:"epoch_length"`
	AverageBlockFrequencyInSeconds float64 `yaml:"average_block_frequency_in_seconds"`
}

func mustLoadConfiguration(configuarionPath string) Configuration {
	bytes, err := os.ReadFile(configuarionPath)

	if err != nil {
		panic(err)
	}

	var configuration Configuration
	err = yaml.Unmarshal(bytes, &configuration)

	if err != nil {
		panic(err)
	}

	return configuration
}
