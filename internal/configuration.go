package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	NodeCount                      int     `yaml:"node_count"`
	SimulationTime                 float64 `yaml:"simulation_time_in_seconds"`
	AverageNetworkLatencyInSeconds int     `yaml:"average_network_latency_in_seconds"`
	AverageBlockFrequencyInSeconds int     `yaml:"average_block_frequency_in_seconds"`
	ChainReogranizationThreshold   int     `yaml:"chain_reorganization_threshold"`

	ProofOfBurnMaturityIntervalInSeconds int `yaml:"proof_of_burn_maturity_interval_in_seconds"`
	ProofOfBurnDecayTimeInSeconds        int `yaml:"proof_of_burn_decay_time_in_seconds"`
}

func LoadConfiguration(configuarionPath string) (Configuration, error) {
	bytes, err := os.ReadFile(configuarionPath)

	if err != nil {
		return Configuration{}, err
	}

	var configuration Configuration
	err = yaml.Unmarshal(bytes, &configuration)

	if err != nil {
		return Configuration{}, err
	}

	return configuration, nil
}
