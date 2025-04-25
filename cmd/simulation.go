package cmd

import (
	"github.com/MarkCDavid/burner/internal"
	"github.com/spf13/cobra"
)

var simulationCmd = &cobra.Command{
	Use: "simulation",
	RunE: func(cmd *cobra.Command, args []string) error {
		simulation := internal.NewSimulation("./configuration/default.yaml")
		internal.TimeFunction(simulation.Simulate, "simulation.Simulate()")
		internal.PrintMemoryUsage()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(simulationCmd)
}
