package cmd

import (
	// "flag"
	// "fmt"
	"os"
	"path/filepath"

	// "path/filepath"

	// "github.com/MarkCDavid/burner/internal/imp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getExecutableName() string {
	return filepath.Base(os.Args[0])
}

var rootCmd = &cobra.Command{
	Use: getExecutableName(),
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func setupLogrus() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
		DisableQuote:  true,
	})
}

func init() {
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		setupLogrus()
	}
}

// const (
// 	FLAG_CONFIGURATION_PATH = "configuration"
// 	FLAG_SEED               = "seed"
// 	FLAG_RUNS               = "runs"
// )

/* func main() { */
// configuration_path := flag.String(FLAG_CONFIGURATION_PATH, "./configuration/default.yaml", "Path to the yaml configuration for the simulator")
// seed := flag.Int64(FLAG_SEED, 0, "Seed for the simulation")
// runs := flag.Int64(FLAG_RUNS, 1, "Run count")
// flag.Parse()
//
// fmt.Println(seed, runs)
//
// if !flagProvided(FLAG_SEED) {
// 	seed = nil
// } else {
// 	_runs := int64(1)
// 	runs = &_runs
// }
//
// for i := int64(0); i < *runs; i++ {
// 	imp.Run(seed)
// }
//
// absolute_configuration_path, err := filepath.Abs(*configuration_path)
// if err != nil {
// 	panic(1)
// }

// for i := int64(0); i < *runs; i++ {
// 	// internal.Simulate(absolute_configuration_path, seed)
// 	internal.PoB(seed)
// }
// }
//
// func flagProvided(name string) bool {
// 	found := false
// 	flag.Visit(func(f *flag.Flag) {
// 		if f.Name == name {
// 			found = true
// 		}
// 	})
// 	return found
// }
