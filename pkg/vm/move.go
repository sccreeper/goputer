package vm

import (
	"encoding/binary"
	c "sccreeper/goputer/pkg/constants"
)

func (m *VM) move() {

	if m.IsImmediate {
		
		if m.RightArg == uint16(c.RVideoText) {
				binary.LittleEndian.PutUint32(m.TextBuffer[:4], m.LeftArgVal)
		} else if m.RightArg == uint16(c.RData) {
				binary.LittleEndian.PutUint32(m.DataBuffer[:4], m.LeftArgVal)
		} else {
			m.Registers[m.RightArg] = m.LeftArgVal
		}

	} else {

		//Copying from buffer -> buffer

		if m.LeftArg == uint16(c.RData) && m.RightArg == uint16(c.RVideoText) && !m.IsImmediate {

			copy(m.TextBuffer[:m.Registers[c.RDataLength]], m.DataBuffer[:m.Registers[c.RDataLength]])

		} else if m.LeftArg == uint16(c.RVideoText) && m.RightArg == uint16(c.RData) && !m.IsImmediate {

			copy(m.DataBuffer[:m.Registers[c.RDataLength]], m.TextBuffer[:m.Registers[c.RDataLength]])

			//Copying from buffer -> register
		} else if m.LeftArg == uint16(c.RData) || m.LeftArg == uint16(c.RVideoText) {
			switch m.LeftArg {
			case uint16(c.RData):
				m.Registers[m.RightArg] = binary.LittleEndian.Uint32(m.DataBuffer[:4])
			case uint16(c.RVideoText):
				m.Registers[m.RightArg] = binary.LittleEndian.Uint32(m.TextBuffer[:4])
			}
			//Copying from register -> buffer
		} else if m.RightArg == uint16(c.RData) || m.RightArg == uint16(c.RVideoText) {
			switch m.RightArg {
			case uint16(c.RData):
				binary.LittleEndian.PutUint32(m.DataBuffer[:4], m.Registers[m.LeftArg])
			case uint16(c.RVideoText):
				binary.LittleEndian.PutUint32(m.TextBuffer[:4], m.Registers[m.LeftArg])
			}

		} else {
			m.Registers[m.RightArg] = m.Registers[m.LeftArg]
		}
	}
}
