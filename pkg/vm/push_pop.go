package vm

import (
	"encoding/binary"
	"errors"
	"sccreeper/govm/pkg/compiler"
	"sccreeper/govm/pkg/constants"
)

func (m *VM) pop_stack() {
	m.Registers[m.ArgSmall0] =
		binary.LittleEndian.Uint32(
			m.MemArray[m.Registers[constants.RStackPointer] : m.Registers[constants.RStackPointer]+4],
		)

	if !(int32(m.Registers[constants.RStackPointer])-4 < 0) {
		m.Registers[constants.RStackPointer] -= 4
	}
}

func (m *VM) push_stack() {

	binary.LittleEndian.PutUint32(m.MemArray[constants.RStackPointer:constants.RStackPointer+4], uint32(m.Registers[m.ArgSmall0]))

	if !(m.Registers[constants.RStackPointer]+4 > compiler.StackSize) {
		m.Registers[constants.RStackPointer] += 4
	} else {
		panic(errors.New("stack overflow"))
	}

}
