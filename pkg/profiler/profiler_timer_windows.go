//go:build windows

package profiler

/*
#include <time_windows.h>
*/
import "C"

var freq uint64

func init() {
	freq = uint64(C.get_perf_freq())
}

func profilerGetTime() uint64 {
	return (uint64(C.get_perf_counter()) * 1e9) / freq
}
