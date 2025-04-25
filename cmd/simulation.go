package cmd

import (
	"github.com/MarkCDavid/burner/internal"
	"github.com/spf13/cobra"
)

var simulationCmd = &cobra.Command{
	Use: "simulation",
	RunE: func(cmd *cobra.Command, args []string) error {
		// simulation := internal.NewSimulation("./configuration/default.yaml")
		// internal.TimeFunction(simulation.Simulate, "simulation.Simulate()")
		// internal.PrintMemoryUsage()

		bitcoin := internal.NewSimulation("./configuration/bitcoin.yaml")
		internal.TimeFunction(bitcoin.Simulate, "bitcoin.Simulate()")
		internal.PrintMemoryUsage()

		// slimcoin := internal.NewSimulation("./configuration/slimcoin.yaml")
		// internal.TimeFunction(slimcoin.Simulate, "slimcoin.Simulate()")
		// internal.PrintMemoryUsage()

		razer := internal.NewSimulation("./configuration/razer.yaml")
		internal.TimeFunction(razer.Simulate, "razer.Simulate()")
		internal.PrintMemoryUsage()

		// solo_razer := internal.NewSimulation("./configuration/solo_razer.yaml")
		// internal.TimeFunction(solo_razer.Simulate, "solo_razer.Simulate()")
		// internal.PrintMemoryUsage()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(simulationCmd)
}
