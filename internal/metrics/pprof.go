package metrics

import (
	"runtime"
	"runtime/pprof"
)

// MemStats returns current memory statistics.
func MemStats() runtime.MemStats {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	return stats
}

// GoroutineCount returns the number of goroutines.
func GoroutineCount() int {
	return runtime.NumGoroutine()
}

// ProfileNames returns available pprof profile names.
func ProfileNames() []string {
	profiles := pprof.Profiles()
	names := make([]string, len(profiles))
	for i, p := range profiles {
		names[i] = p.Name()
	}
	return names
}
