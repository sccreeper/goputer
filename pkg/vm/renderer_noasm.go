//go:build !amd64

package vm

const haveArchVideoClear = false
const haveArchVideoArea = false

func archVideoClear(array *byte, r uint8, g uint8, b uint8) {
	panic("not implemented")
}

func archVideoArea(array *byte, red uint8, green uint8, blue uint8, alpha uint8, x uint32, y uint32, x1 uint32, y1 uint32) {
	panic("no implemented")
}
