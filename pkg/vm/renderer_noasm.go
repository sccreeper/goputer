//go:build !amd64

package vm

const haveArchVideoClear = false

func archVideoClear(array *byte, r uint8, g uint8, b uint8) {
	panic("not implemented")
}
