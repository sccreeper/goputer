//go:build amd64

package vm

import (
	"fmt"
	"os"
	"sccreeper/goputer/pkg/vm/asm"

	"golang.org/x/sys/cpu"
)

const haveArchVideoClear = true
const haveArchVideoArea = true

func init() {
	if !(cpu.X86.HasSSE2 && cpu.X86.HasAVX && cpu.X86.HasAVX2) {
		fmt.Println("CPU does not support SSE2, AVX, or AVX2.")
		os.Exit(1)
	}
}

func archVideoClear(array *byte, r uint8, g uint8, b uint8) {
	asm.VideoClearAsm(array, r, g, b)
}

func archVideoAreaNoAlpha(array *byte, red uint8, green uint8, blue uint8, x uint32, y uint32, x1 uint32, y1 uint32) {
	asm.VideoAreaAsm(array, red, green, blue, x, y, x1, y1)
}
