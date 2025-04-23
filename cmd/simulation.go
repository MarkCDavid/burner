package cmd

import (
	"time"

	"github.com/MarkCDavid/burner/internal"
	// "github.com/MarkCDavid/burner/internal_bkp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var simulationCmd = &cobra.Command{
	Use: "simulation",
	RunE: func(cmd *cobra.Command, args []string) error {
		// seed := int64(1745405314273631501)

		// start := time.Now()
		// internal_bkp.Simulate("./configuration/default.yaml", nil)
		// elapsed := time.Since(start)
		// logrus.Warnf("Old: %s", elapsed)
		// internal.PrintMemUsage()

		start := time.Now()
		internal.Simulate("./configuration/default.yaml", nil)
		// internal.Simulate("./configuration/default.yaml", &seed)
		elapsed := time.Since(start)
		logrus.Warnf("New: %s", elapsed)
		internal.PrintMemUsage()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(simulationCmd)
}
