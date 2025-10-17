package vm

import (
	"sccreeper/goputer/pkg/compiler"
	c "sccreeper/goputer/pkg/constants"
)

func (m *VM) call(addr uint32, dest uint32) {
	m.pushCall(addr)
	m.Registers[c.RProgramCounter] = dest
}

func (m *VM) conditionalCall() bool {

	if m.Registers[c.RAccumulator] != 0 {

		var addressVal uint32

		if m.IsImmediate {
			addressVal = m.LongArgVal
		} else if uint16(m.LongArg) < MaxRegister {
			addressVal = m.Registers[m.LongArg]
		} else {
			addressVal = m.LongArg
		}

		m.pushCall(m.Registers[c.RProgramCounter] + compiler.InstructionLength)
		m.Registers[c.RProgramCounter] = addressVal

		return true
	} else {
		return false
	}
}

// Jumps
func (m *VM) jump() {

	var addressVal uint32

	if m.IsImmediate {
		addressVal = m.LongArgVal
	} else if uint16(m.LongArg) < MaxRegister {
		addressVal = m.Registers[m.LongArg]
	} else {
		addressVal = m.LongArg
	}

	m.Registers[c.RProgramCounter] = addressVal
}

func (m *VM) conditionalJump() bool {

	if m.Registers[c.RAccumulator] != 0 {
		var addressVal uint32

		if m.IsImmediate {
			addressVal = m.LongArgVal
		} else if uint16(m.LongArg) < MaxRegister {
			addressVal = m.Registers[m.LongArg]
		} else {
			addressVal = m.LongArg
		}

		m.Registers[c.RProgramCounter] = addressVal
		return true
	} else {
		return false
	}

}
