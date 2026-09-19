package cmd

import (
	"github.com/MarkCDavid/burner/internal"
	"github.com/spf13/cobra"
)

var seedFlag int64

var simulationCmd = &cobra.Command{
	Use: "simulation",
	RunE: func(cmd *cobra.Command, args []string) error {
		simulation := internal.NewSimulation(args[0], seedFlag)
		internal.TimeFunction(simulation.Simulate, args[0])
		internal.PrintMemoryUsage()

		return nil
	},
}

func init() {
	simulationCmd.Flags().Int64Var(&seedFlag, "seed", 0, "Override the seed from the configuration file")
	rootCmd.AddCommand(simulationCmd)
}
