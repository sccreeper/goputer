package vm

import (
	c "sccreeper/goputer/pkg/constants"
)

//Shifting left and shifting right

func (m *VM) shiftLeft() {
	switch m.ArgSmall0 {
	case uint16(c.RData):
		copy(m.DataBuffer[:], append(m.DataBuffer[m.Registers[m.ArgSmall1]:], make([]byte, 128-len(m.DataBuffer[m.Registers[m.ArgSmall1]:]))...))
	case uint16(c.RVideoText):
		copy(m.TextBuffer[:], append(m.TextBuffer[m.Registers[m.ArgSmall1]:], make([]byte, 128-len(m.TextBuffer[m.Registers[m.ArgSmall1]:]))...))
	default:
		m.Registers[c.RAccumulator] = m.Registers[m.ArgSmall0] << m.Registers[m.ArgSmall1]
	}

}

func (m *VM) shiftRight() {
	switch m.ArgSmall0 {
	case uint16(c.RData):
		copy(m.DataBuffer[:], append(make([]byte, m.Registers[m.ArgSmall1]), m.DataBuffer[:128-m.Registers[m.Registers[m.ArgSmall1]]]...))
	case uint16(c.RVideoText):
		copy(m.TextBuffer[:], append(make([]byte, m.Registers[m.ArgSmall1]), m.TextBuffer[:128-m.Registers[m.Registers[m.ArgSmall1]]]...))
	default:
		m.Registers[c.RAccumulator] = m.Registers[m.ArgSmall0] >> m.Registers[m.ArgSmall1]
	}
}
