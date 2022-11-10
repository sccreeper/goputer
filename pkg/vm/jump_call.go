package vm

import (
	"encoding/binary"
	"sccreeper/goputer/pkg/compiler"
	c "sccreeper/goputer/pkg/constants"
)

// Calls
func (m *VM) call() {
	m.Registers[c.RCallStackPointer] += 4
	binary.LittleEndian.PutUint32(
		m.MemArray[m.Registers[c.RCallStackPointer]:m.Registers[c.RCallStackPointer]+4],
		m.Registers[c.RProgramCounter]+compiler.InstructionLength)
	m.Registers[c.RProgramCounter] = m.ArgLarge
}

func (m *VM) conditional_call() {
	if m.Registers[c.RAccumulator] != 0 {
		m.Registers[c.RCallStackPointer] += 4
		binary.LittleEndian.PutUint32(
			m.MemArray[m.Registers[c.RCallStackPointer]:m.Registers[c.RCallStackPointer]+4],
			m.Registers[c.RProgramCounter]+compiler.InstructionLength)
		m.Registers[c.RProgramCounter] = m.ArgLarge
	}
}

// Jumps
func (m *VM) jump() {
	m.Registers[c.RProgramCounter] = m.ArgLarge
}

func (m *VM) conditional_jump() {

	if m.Registers[c.RAccumulator] != 0 {
		m.Registers[c.RProgramCounter] = m.ArgLarge
	}

}
