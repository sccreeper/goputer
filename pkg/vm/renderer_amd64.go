//go:build amd64

package vm

import "sccreeper/goputer/pkg/vm/asm"

const haveArchVideoClear = true

func archVideoClear(array *byte, r uint8, g uint8, b uint8) {
	asm.VideoClearAsm(array, r, g, b)
}
