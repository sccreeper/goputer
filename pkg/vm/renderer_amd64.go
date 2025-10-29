//go:build amd64

package vm

import (
	"fmt"
	"os"
	"sccreeper/goputer/pkg/vm/asm"

	"golang.org/x/sys/cpu"
)

const haveArchVideoClear = true

func init() {
	if !(cpu.X86.HasSSE2 && cpu.X86.HasAVX && cpu.X86.HasAVX2) {
		fmt.Println("CPU does not support SSE2, AVX, or AVX2.")
		os.Exit(1)
	}
}

func archVideoClear(array *byte, r uint8, g uint8, b uint8) {
	asm.VideoClearAsm(array, r, g, b)
}
