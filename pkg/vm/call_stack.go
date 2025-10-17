package vm

import (
	"encoding/binary"
	"sccreeper/goputer/pkg/compiler"
	c "sccreeper/goputer/pkg/constants"
)

func (m *VM) popCall() {

	if m.Registers[c.RCallStackPointer]-4 < uint32(c.RCallStackZeroPointer) {
		m.Registers[c.RCallStackPointer] = uint32(c.RCallStackZeroPointer)
	} else {
		m.Registers[c.RCallStackPointer] -= 4
	}

	m.Registers[c.RProgramCounter] = binary.LittleEndian.Uint32(
		m.MemArray[m.Registers[c.RCallStackPointer] : m.Registers[c.RCallStackPointer]+4],
	)

	copy(m.MemArray[m.Registers[c.RCallStackPointer]:m.Registers[c.RCallStackPointer]+4], []byte{0, 0, 0, 0})

}

func (m *VM) pushCall(addr uint32) {

	if m.Registers[c.RCallStackPointer]+4 > compiler.DataStackSize+compiler.CallStackSize+VideoBufferSize {
		panic("call stack overflow")
	}

	binary.LittleEndian.PutUint32(
		m.MemArray[m.Registers[c.RCallStackPointer]:m.Registers[c.RCallStackPointer]+4],
		addr,
	)

	m.Registers[c.RCallStackPointer] += 4

}
