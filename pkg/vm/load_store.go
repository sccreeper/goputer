package vm

import (
	"encoding/binary"
	"sccreeper/goputer/pkg/constants"
)

func (m *VM) load() {

	data_length := binary.LittleEndian.Uint32(m.MemArray[m.ArgLarge : m.ArgLarge+4])
	m.Registers[constants.RDataLength] = data_length

	copy(m.DataBuffer[:data_length], m.MemArray[m.ArgLarge+4:m.ArgLarge+data_length+4])

}

func (m *VM) store() {

	data_length := binary.LittleEndian.Uint32(m.MemArray[m.ArgLarge : m.ArgLarge+4])
	m.Registers[constants.RDataLength] = data_length

	copy(m.MemArray[m.ArgLarge+4:m.ArgLarge+data_length+4], m.DataBuffer[:data_length])

}
