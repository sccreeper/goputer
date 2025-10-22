package vm

import (
	"encoding/binary"
	"sccreeper/goputer/pkg/compiler"
	c "sccreeper/goputer/pkg/constants"
)

func (m *VM) popStack() {
	if m.Registers[c.RStackPointer]-4 >= m.Registers[c.RStackZeroPointer] {
		m.Registers[m.LeftArg] =
			binary.LittleEndian.Uint32(
				m.MemArray[m.Registers[c.RStackPointer]-4 : m.Registers[c.RStackPointer]],
			)

		m.Registers[c.RStackPointer] -= 4
	} else {
		m.Registers[c.RStackPointer] = m.Registers[c.RStackZeroPointer] // prevent underflow
	}

}

func (m *VM) pushStack() {

	if m.Registers[c.RStackPointer]+4 < compiler.MemOffset {
		binary.LittleEndian.PutUint32(
			m.MemArray[m.Registers[c.RStackPointer]:m.Registers[c.RStackPointer]+4],
			m.LeftArgVal,
		)

		m.Registers[c.RStackPointer] += 4
	} else {
		panic("data stack overflow")
	}

}
