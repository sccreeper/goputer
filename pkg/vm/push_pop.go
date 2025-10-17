package vm

import (
	"encoding/binary"
	"sccreeper/goputer/pkg/compiler"
	c "sccreeper/goputer/pkg/constants"
)

func (m *VM) popStack() {
	if !(m.Registers[c.RStackPointer]-4 < m.Registers[c.RStackZeroPointer]) {
		m.Registers[c.RStackPointer] -= 4
	} else {
		m.Registers[c.RStackPointer] = m.Registers[c.RStackZeroPointer] // prevent underflow
	}

	m.Registers[m.LeftArg] =
		binary.LittleEndian.Uint32(
			m.MemArray[m.Registers[c.RStackPointer]-4 : m.Registers[c.RStackPointer]],
		)

}

func (m *VM) pushStack() {

	binary.LittleEndian.PutUint32(
		m.MemArray[m.Registers[c.RStackPointer]:m.Registers[c.RStackPointer]+4],
		uint32(m.LeftArgVal),
	)

	if m.Registers[c.RStackPointer]+4 < compiler.DataStackSize+VideoBufferSize {
		m.Registers[c.RStackPointer] += 4
	} else {
		panic("data stack overflow")
	}

}
