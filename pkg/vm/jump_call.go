package vm

import (
	"encoding/binary"
	"sccreeper/goputer/pkg/compiler"
	c "sccreeper/goputer/pkg/constants"
)

// Calls the address of m.ArgLarge
func (m *VM) call() {

	var addressVal uint32

	if uint16(m.ArgLarge) < RegisterCount {
		addressVal = m.Registers[m.ArgLarge]
	} else {
		addressVal = m.ArgLarge
	}

	var increment uint32

	if m.Opcode != c.ICall {
		increment = 0
	} else {
		increment = compiler.InstructionLength
	}

	m.Registers[c.RCallStackPointer] += 4

	binary.LittleEndian.PutUint32(
		m.MemArray[m.Registers[c.RCallStackPointer]:m.Registers[c.RCallStackPointer]+4],
		m.Registers[c.RProgramCounter]+increment)

	m.Registers[c.RProgramCounter] = addressVal
}

func (m *VM) conditionalCall() bool {

	if m.Registers[c.RAccumulator] != 0 {

		var addressVal uint32

		if uint16(m.ArgLarge) < RegisterCount {
			addressVal = m.Registers[m.ArgLarge]
		} else {
			addressVal = m.ArgLarge
		}

		m.Registers[c.RCallStackPointer] += 4

		binary.LittleEndian.PutUint32(
			m.MemArray[m.Registers[c.RCallStackPointer]:m.Registers[c.RCallStackPointer]+4],
			m.Registers[c.RProgramCounter]+compiler.InstructionLength)
		m.Registers[c.RProgramCounter] = addressVal

		return true
	} else {
		return false
	}
}

// Jumps
func (m *VM) jump() {

	var addressVal uint32

	if uint16(m.ArgLarge) < RegisterCount {
		addressVal = m.Registers[m.ArgLarge]
	} else {
		addressVal = m.ArgLarge
	}

	m.HandlingInterrupt = false

	m.Registers[c.RProgramCounter] = addressVal
}

func (m *VM) conditionalJump() bool {

	if m.Registers[c.RAccumulator] != 0 {
		var addressVal uint32

		if uint16(m.ArgLarge) < RegisterCount {
			addressVal = m.Registers[m.ArgLarge]
		} else {
			addressVal = m.ArgLarge
		}

		m.HandlingInterrupt = false
		m.Registers[c.RProgramCounter] = addressVal
		return true
	} else {
		return false
	}

}
