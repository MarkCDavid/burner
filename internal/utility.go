package internal

import (
	"fmt"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

func TimeFunction(function func(), name string) {
	start := time.Now()
	function()
	elapsed := time.Since(start)

	logrus.Infof("%s ran for %s", name, elapsed)

}

func PrintMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logrus.Infof("Memory Usage: %s | GC ran %v times", FormatBytes(m.Sys), m.NumGC)
}

const KiB = 1024

func FormatBytes(size uint64) string {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
	i := 0

	for size >= KiB && i < len(units)-1 {
		size /= KiB
		i++
	}

	return fmt.Sprintf("%v %s", size, units[i])
}
