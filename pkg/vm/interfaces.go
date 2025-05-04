package vm

import "io"

// Implements several std/io interfaces

func (m *VM) ReadAt(p []byte, off int64) (n int, err error) {

	if off > int64(MemSize) {
		return 0, io.EOF
	}

	var limit int64

	if len(p) > int(MemSize)-int(off) {
		limit = int64(MemSize)
	} else {
		limit = int64(len(p))
	}

	var bytesWritten int = copy(
		p[:],
		m.MemArray[off:off+limit],
	)

	return bytesWritten, nil

}
