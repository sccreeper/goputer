package vm

import (
	c "sccreeper/goputer/pkg/constants"
)

//Shifting left and shifting right

func (m *VM) shiftLeft() {
	switch m.LeftArg {
	case uint16(c.RData):
		copy(m.DataBuffer[:], append(m.DataBuffer[m.RightArgVal:], make([]byte, 128-len(m.DataBuffer[m.RightArgVal:]))...))
	case uint16(c.RVideoText):
		copy(m.TextBuffer[:], append(m.TextBuffer[m.RightArgVal:], make([]byte, 128-len(m.TextBuffer[m.RightArgVal:]))...))
	default:
		m.Registers[c.RAccumulator] = m.LeftArgVal << m.RightArgVal
	}

}

func (m *VM) shiftRight() {
	switch m.LeftArg {
	case uint16(c.RData):
		copy(m.DataBuffer[:], append(make([]byte, m.RightArgVal), m.DataBuffer[:128-m.RightArgVal]...))
	case uint16(c.RVideoText):
		copy(m.TextBuffer[:], append(make([]byte, m.RightArgVal), m.TextBuffer[:128-m.RightArgVal]...))
	default:
		m.Registers[c.RAccumulator] = m.LeftArgVal >> m.RightArgVal
	}
}
