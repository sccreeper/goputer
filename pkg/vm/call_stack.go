package vm

import (
	"encoding/binary"
	c "sccreeper/goputer/pkg/constants"
)

// Set program pointer to previous position and then pop from stack.
func (m *VM) popCall() {

	//Get address at call stack pointer
	m.Registers[c.RProgramCounter] = binary.LittleEndian.Uint32(
		m.MemArray[m.Registers[c.RCallStackPointer] : m.Registers[c.RCallStackPointer]+4],
	)

	//Overwrite call stack pointer
	copy(m.MemArray[m.Registers[c.RCallStackPointer]:m.Registers[c.RCallStackPointer]+4], []byte{0, 0, 0, 0})

	//Finally decrement the call stack pointer
	if m.Registers[c.RCallStackPointer]-4 < uint32(c.RCallStackZeroPointer) {
		m.Registers[c.RCallStackPointer] = uint32(c.RCallStackPointer)
	} else {
		m.Registers[c.RCallStackPointer] -= 4
	}

}
