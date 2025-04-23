package cmd

import (
	"github.com/MarkCDavid/burner/internal"
	"github.com/spf13/cobra"
)

var simulationCmd = &cobra.Command{
	Use: "simulation",
	RunE: func(cmd *cobra.Command, args []string) error {
		// seed := int64(1745405314273631501)

		internal.Simulate("./configuration/default.yaml", nil)
		// internal.Simulate("./configuration/default.yaml", &seed)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(simulationCmd)
}
