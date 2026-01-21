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
const haveArchVideoAreaAlpha = true

func init() {
	if !(cpu.X86.HasSSE2 && cpu.X86.HasAVX && cpu.X86.HasAVX2 && cpu.X86.HasSSE41) {
		fmt.Println("x86 CPU must support SSE2, SSE4.1, AVX, and AVX2 instruction extensions.")
		os.Exit(1)
	}
}

func archVideoClear(array *byte, r uint8, g uint8, b uint8) {
	asm.VideoClearAsm(array, r, g, b)
}

func archVideoAreaNoAlpha(array *byte, red uint8, green uint8, blue uint8, x uint32, y uint32, x1 uint32, y1 uint32) {
	asm.VideoAreaAsm(array, red, green, blue, x, y, x1, y1)
}

func archVideoAreaAlpha(array *byte, red uint8, green uint8, blue uint8, alpha uint8, x uint32, y uint32, x1 uint32, y1 uint32) {
	asm.VideoAreaAlphaAsm(array, red, green, blue, alpha, x, y, x1, y1)
}
