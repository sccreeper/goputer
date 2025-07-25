package vm

import (
	"encoding/binary"
	c "sccreeper/goputer/pkg/constants"
)

// Load r00 bytes from address 64
//lda $64 r00
// Load 64 bytes from r00
//lda r00 $64

func (m *VM) load() {

	// This means we are unable to "directly address" the first 57 bytes of memory (register address space).
	// Addressing static stack value
	if m.LeftArg > MaxRegister && !m.IsImmediate {
		dataLength := binary.LittleEndian.Uint32(m.MemArray[m.LongArg : m.LongArg+4])
		m.Registers[c.RDataLength] = dataLength

		if dataLength > 128 {
			copy(m.DataBuffer[:], m.MemArray[m.LongArg+4:m.LongArg+128+4])
		} else {
			copy(m.DataBuffer[:dataLength], m.MemArray[m.LongArg+4:m.LongArg+dataLength+4])
		}

		m.Registers[c.RDataPointer] = m.LongArg + 4 // If we want data length, it is in dl
	} else { // Addressing main memory
		copy(m.DataBuffer[:m.RightArgVal], m.MemArray[m.LeftArgVal:m.LeftArgVal+m.RightArgVal])

		m.Registers[c.RDataPointer] = m.LeftArgVal
		m.Registers[c.RDataLength] = m.RightArgVal
	}

}

func (m *VM) store() {

	// Addressing static stack
	if m.LeftArg > MaxRegister && !m.IsImmediate {
		dataLength := binary.LittleEndian.Uint32(m.MemArray[m.LongArg : m.LongArg+4])
		m.Registers[c.RDataLength] = dataLength

		if dataLength > 128 {
			dataLength = 128
		}

		copy(m.MemArray[m.LongArg+4:m.LongArg+dataLength+4], m.DataBuffer[:dataLength])

		m.Registers[c.RDataPointer] = m.LongArg
	} else { // Addressing main memory

		copy(m.MemArray[m.LeftArgVal:m.LeftArgVal+m.RightArgVal], m.DataBuffer[:m.RightArgVal])

		m.Registers[c.RDataPointer] = m.LeftArgVal
		m.Registers[c.RDataLength] = m.RightArgVal
	}

}
