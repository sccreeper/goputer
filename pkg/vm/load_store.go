package vm

import (
	"encoding/binary"
	c "sccreeper/goputer/pkg/constants"
)

func (m *VM) load() {

	if m.ArgLarge > uint32(RegisterCount) {
		data_length := binary.LittleEndian.Uint32(m.MemArray[m.ArgLarge : m.ArgLarge+4])
		m.Registers[c.RDataLength] = data_length

		copy(m.DataBuffer[:data_length], m.MemArray[m.ArgLarge+4:m.ArgLarge+data_length+4])

		m.Registers[c.RDataPointer] = m.ArgLarge
	} else {
		copy(m.DataBuffer[:m.Registers[m.ArgLarge]], m.MemArray[m.Registers[c.RDataPointer]:m.Registers[c.RDataPointer]+m.Registers[m.ArgLarge]])
	}

}

func (m *VM) store() {

	if m.ArgLarge > uint32(RegisterCount) {
		data_length := binary.LittleEndian.Uint32(m.MemArray[m.ArgLarge : m.ArgLarge+4])
		m.Registers[c.RDataLength] = data_length

		copy(m.MemArray[m.ArgLarge+4:m.ArgLarge+data_length+4], m.DataBuffer[:data_length])

		m.Registers[c.RDataPointer] = m.ArgLarge
	} else {
		copy(m.MemArray[m.Registers[c.RDataPointer]:m.Registers[c.RDataPointer]+m.Registers[m.ArgLarge]], m.DataBuffer[:m.Registers[m.ArgLarge]])
	}

}
