package cmd

import (
	"github.com/MarkCDavid/burner/internal"
	"github.com/spf13/cobra"
)

var simulationCmd = &cobra.Command{
	Use: "simulation",
	RunE: func(cmd *cobra.Command, args []string) error {
		internal.Simulate("./configuration/default.yaml", nil)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(simulationCmd)
}
