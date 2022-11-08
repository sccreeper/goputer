package vm

import "encoding/binary"

func (m *VM) load() {

	data_length := binary.LittleEndian.Uint32(m.MemArray[m.ArgLarge : m.ArgLarge+4])

	copy(m.DataBuffer[:], m.MemArray[m.ArgLarge+4:m.ArgLarge+data_length+4])

}

func (m *VM) store() {

	data_length := binary.LittleEndian.Uint32(m.MemArray[m.ArgLarge : m.ArgLarge+4])

	copy(m.MemArray[m.ArgLarge+4:m.ArgLarge+data_length+4], m.DataBuffer[:data_length])

}
