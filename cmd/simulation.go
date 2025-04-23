package cmd

import (
	"time"

	"github.com/MarkCDavid/burner/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var simulationCmd = &cobra.Command{
	Use: "simulation",
	RunE: func(cmd *cobra.Command, args []string) error {
		// seed := int64(1745405314273631501)

		start := time.Now()
		internal.Simulate("./configuration/default.yaml", nil)
		// internal.Simulate("./configuration/default.yaml", &seed)
		elapsed := time.Since(start)
		logrus.Warnf("Took %s", elapsed)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(simulationCmd)
}
