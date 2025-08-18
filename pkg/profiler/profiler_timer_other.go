//go:build !windows

package profiler

import (
	"time"
)

func profilerGetTime() uint64 {
	return uint64(time.Now().UnixNano())
}
