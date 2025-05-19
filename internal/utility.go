package internal

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

func TimeFunction(function func(), name string) {
	logrus.Info()
	logrus.Info()
	logrus.Infof("Running %s", name)
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

func Filter[T any](slice []T, keep func(T) bool) []T {
	n := 0
	for _, v := range slice {
		if keep(v) {
			slice[n] = v
			n++
		}
	}
	return slice[:n]
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

func ClampPositiveFloat64(x float64) float64 {
	switch {
	case x == 0:
		return math.SmallestNonzeroFloat64
	case math.IsInf(x, 1):
		return math.MaxFloat64
	case math.IsInf(x, -1):
		return -math.MaxFloat64
	default:
		return x
	}
}

func ClampFloat64(x float64, min float64, max float64) float64 {
	switch {
	case x < min:
		return min
	case x > max:
		return max
	default:
		return x
	}
}
