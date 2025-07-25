package vm

import (
	c "sccreeper/goputer/pkg/constants"
)

//Shifting left and shifting right

func (m *VM) shiftLeft() {
	switch m.LeftArg {
	case uint16(c.RData):
		copy(m.DataBuffer[:], append(m.DataBuffer[m.Registers[m.RightArg]:], make([]byte, 128-len(m.DataBuffer[m.Registers[m.RightArg]:]))...))
	case uint16(c.RVideoText):
		copy(m.TextBuffer[:], append(m.TextBuffer[m.Registers[m.RightArg]:], make([]byte, 128-len(m.TextBuffer[m.Registers[m.RightArg]:]))...))
	default:
		m.Registers[c.RAccumulator] = m.Registers[m.LeftArg] << m.Registers[m.RightArg]
	}

}

func (m *VM) shiftRight() {
	switch m.LeftArg {
	case uint16(c.RData):
		copy(m.DataBuffer[:], append(make([]byte, m.Registers[m.RightArg]), m.DataBuffer[:128-m.Registers[m.Registers[m.RightArg]]]...))
	case uint16(c.RVideoText):
		copy(m.TextBuffer[:], append(make([]byte, m.Registers[m.RightArg]), m.TextBuffer[:128-m.Registers[m.Registers[m.RightArg]]]...))
	default:
		m.Registers[c.RAccumulator] = m.Registers[m.LeftArg] >> m.Registers[m.RightArg]
	}
}
