package vm

import (
	"encoding/binary"
	c "sccreeper/goputer/pkg/constants"
)

func (m *VM) load() {

	// This means we are unable to "directly address" the first 57 bytes of memory (register address space).
	// Addressing stack
	if m.ArgSmall0 > RegisterCount {
		dataLength := binary.LittleEndian.Uint32(m.MemArray[m.ArgLarge : m.ArgLarge+4])
		m.Registers[c.RDataLength] = dataLength

		if dataLength > 128 {
			copy(m.DataBuffer[:], m.MemArray[m.ArgLarge+4:m.ArgLarge+128+4])
		} else {
			copy(m.DataBuffer[:dataLength], m.MemArray[m.ArgLarge+4:m.ArgLarge+dataLength+4])
		}

		m.Registers[c.RDataPointer] = m.ArgLarge + 4 // If we want data length, it is in dl
	} else { // Addressing main memory
		copy(m.DataBuffer[:m.Registers[m.ArgSmall1]], m.MemArray[m.Registers[m.ArgSmall0]:m.Registers[m.ArgSmall0]+m.Registers[m.ArgSmall1]])

		m.Registers[c.RDataPointer] = m.Registers[m.ArgSmall0]
		m.Registers[c.RDataLength] = m.Registers[m.ArgSmall1]
	}

}

func (m *VM) store() {

	// Addressing stack
	if m.ArgSmall0 > RegisterCount {
		dataLength := binary.LittleEndian.Uint32(m.MemArray[m.ArgLarge : m.ArgLarge+4])
		m.Registers[c.RDataLength] = dataLength

		if dataLength > 128 {
			dataLength = 128
		}

		copy(m.MemArray[m.ArgLarge+4:m.ArgLarge+dataLength+4], m.DataBuffer[:dataLength])

		m.Registers[c.RDataPointer] = m.ArgLarge
	} else { // Addressing main memory
		copy(m.MemArray[m.Registers[m.ArgSmall0]:m.Registers[m.ArgSmall0]+m.Registers[m.ArgSmall1]], m.DataBuffer[:m.Registers[m.ArgSmall1]])

		m.Registers[c.RDataPointer] = m.Registers[m.ArgSmall0]
		m.Registers[c.RDataLength] = m.Registers[m.ArgSmall1]
	}

}
