package vm

import (
	"encoding/binary"
	"errors"
	"sccreeper/goputer/pkg/compiler"
	c "sccreeper/goputer/pkg/constants"
)

func (m *VM) popStack() {
	m.Registers[m.LeftArg] =
		binary.LittleEndian.Uint32(
			m.MemArray[m.Registers[c.RStackPointer]-4 : m.Registers[c.RStackPointer]],
		)

	if !(int32(m.Registers[c.RStackPointer])-4 < 0) {
		m.Registers[c.RStackPointer] -= 4
	}

}

func (m *VM) pushStack() {

	binary.LittleEndian.PutUint32(m.MemArray[m.Registers[c.RStackPointer]:m.Registers[c.RStackPointer]+4], uint32(m.LeftArgVal))

	if m.Registers[c.RStackPointer]+4 < compiler.DataStackSize + VideoBufferSize {
		m.Registers[c.RStackPointer] += 4
	} else {
		panic(errors.New("stack overflow"))
	}

}
