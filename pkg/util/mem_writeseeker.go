package util

import (
	"errors"
	"io"
)

// Memory based WriteSeeker for compatibility between gpimg and WASM

type MemWriteSeeker struct {
	buffer []byte
	offset int64
}

func NewMemWriteSeeker() *MemWriteSeeker {

	return &MemWriteSeeker{
		buffer: make([]byte, 0),
		offset: 0,
	}

}

func (m *MemWriteSeeker) Write(p []byte) (n int, err error) {

	newOffset := int(m.offset) + len(p)

	if newOffset > len(m.buffer) {
		newBuf := make([]byte, newOffset)
		copy(newBuf, m.buffer)
		m.buffer = newBuf
	}

	copy(m.buffer[m.offset:], p)

	m.offset += int64(len(p))

	return len(p), nil

}

func (m *MemWriteSeeker) Seek(offset int64, whence int) (int64, error) {

	var newOffset int64

	switch whence {
	case io.SeekStart:
		newOffset = offset
	case io.SeekCurrent:
		newOffset = m.offset + offset
	case io.SeekEnd:
		newOffset = int64(len(m.buffer)) + offset
	default:
		return 0, errors.New("invalid whence")
	}

	if newOffset < 0 {
		return 0, errors.New("negative offset")
	}

	m.offset = newOffset

	return newOffset, nil
}

func (m *MemWriteSeeker) Bytes() []byte {
	return m.buffer
}
