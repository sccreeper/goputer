package vm

import (
	"encoding/binary"
	c "sccreeper/goputer/pkg/constants"
)

func (m *VM) load() {

	if m.ArgLarge > uint32(RegisterCount) {
		dataLength := binary.LittleEndian.Uint32(m.MemArray[m.ArgLarge : m.ArgLarge+4])
		m.Registers[c.RDataLength] = dataLength

		if dataLength > 128 {
			copy(m.DataBuffer[:], m.MemArray[m.ArgLarge+4:m.ArgLarge+128+4])
			m.Registers[c.RDataLength] = 128
		} else {
			copy(m.DataBuffer[:dataLength], m.MemArray[m.ArgLarge+4:m.ArgLarge+dataLength+4])
		}

		m.Registers[c.RDataPointer] = m.ArgLarge
	} else {
		copy(m.DataBuffer[:m.Registers[m.ArgLarge]], m.MemArray[m.Registers[c.RDataPointer]:m.Registers[c.RDataPointer]+m.Registers[m.ArgLarge]])
	}

}

func (m *VM) store() {

	if m.ArgLarge > uint32(RegisterCount) {
		dataLength := binary.LittleEndian.Uint32(m.MemArray[m.ArgLarge : m.ArgLarge+4])
		m.Registers[c.RDataLength] = dataLength

		if dataLength > 128 {
			dataLength = 128
		}

		copy(m.MemArray[m.ArgLarge+4:m.ArgLarge+dataLength+4], m.DataBuffer[:dataLength])

		m.Registers[c.RDataPointer] = m.ArgLarge
	} else {
		copy(m.MemArray[m.Registers[c.RDataPointer]:m.Registers[c.RDataPointer]+m.Registers[m.ArgLarge]], m.DataBuffer[:m.Registers[m.ArgLarge]])
	}

}
