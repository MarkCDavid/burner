package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Seed int64 `yaml:"seed"`

	NodeCount      int64   `yaml:"node_count"`
	SimulationTime float64 `yaml:"simulation_time_in_seconds"`

	ChainReogranizationThreshold   int64   `yaml:"chain_reorganization_threshold"`
	AverageNetworkLatencyInSeconds float64 `yaml:"average_network_latency_in_seconds"`

	AverageTransactionsPerSecond int64 `yaml:"average_transactions_per_second"`
	MaximumTransactionsPerBlock  int64 `yaml:"maximum_transaction_per_block"`

	ProofOfWork         Consensus_PoW_Configuration  `yaml:"proof_of_work"`
	SlimcoinProofOfBurn Consensus_SPoB_Configuration `yaml:"slimcoin_proof_of_burn"`
	RazerProofOfBurn    Consensus_RPoB_Configuration `yaml:"razer_proof_of_burn"`
	PricingProofOfBurn  Consensus_PPoB_Configuration `yaml:"pricing_proof_of_burn"`
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
