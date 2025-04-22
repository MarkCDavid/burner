package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	NodeCount      int     `yaml:"node_count"`
	SimulationTime float64 `yaml:"simulation_time_in_seconds"`

	AverageNetworkLatencyInSeconds int `yaml:"average_network_latency_in_seconds"`
	LowerBoundNetworkLatency       float64
	UpperBoundNetworkLatency       float64
	InverseAverageNetworkLatency   float64

	AverageBlockFrequencyInSeconds int `yaml:"average_block_frequency_in_seconds"`
	InverseAverageBlockFrequency   float64

	ChainReogranizationThreshold int `yaml:"chain_reorganization_threshold"`
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

	configuration.LowerBoundNetworkLatency = float64(configuration.AverageNetworkLatencyInSeconds) * float64(0.5)
	configuration.UpperBoundNetworkLatency = float64(configuration.AverageNetworkLatencyInSeconds) * float64(2.0)
	configuration.InverseAverageNetworkLatency = 1.0 / float64(configuration.AverageNetworkLatencyInSeconds)

	configuration.InverseAverageBlockFrequency = 1.0 / (float64(configuration.AverageBlockFrequencyInSeconds) * float64(configuration.NodeCount))

	return configuration
}
