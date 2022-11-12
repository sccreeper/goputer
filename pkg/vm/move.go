package vm

import (
	"encoding/binary"
	c "sccreeper/goputer/pkg/constants"
)

func (m *VM) move() {

	//Copying from buffer -> buffer

	if m.ArgSmall0 == uint16(c.RData) && m.ArgSmall1 == uint16(c.RVideoText) {

		copy(m.TextBuffer[:m.Registers[c.RDataLength]], m.DataBuffer[:m.Registers[c.RDataLength]])

	} else if m.ArgSmall0 == uint16(c.RVideoText) && m.ArgSmall1 == uint16(c.RData) {

		copy(m.DataBuffer[:m.Registers[c.RDataLength]], m.TextBuffer[:m.Registers[c.RDataLength]])

		//Copying from buffer -> register
	} else if m.ArgSmall0 == uint16(c.RData) || m.ArgSmall0 == uint16(c.RVideoText) {
		switch m.ArgSmall0 {
		case uint16(c.RData):
			m.Registers[m.ArgSmall1] = binary.LittleEndian.Uint32(m.DataBuffer[:4])
		case uint16(c.RVideoText):
			m.Registers[m.ArgSmall1] = binary.LittleEndian.Uint32(m.TextBuffer[:4])
		}
		//Copying from register -> buffer
	} else if m.ArgSmall1 == uint16(c.RData) || m.ArgSmall1 == uint16(c.RVideoText) {
		switch m.ArgSmall1 {
		case uint16(c.RData):
			binary.BigEndian.PutUint32(m.DataBuffer[:4], m.Registers[m.ArgSmall0])
		case uint16(c.RVideoText):
			binary.BigEndian.PutUint32(m.TextBuffer[:4], m.Registers[m.ArgSmall0])
		}

	} else {
		m.Registers[m.ArgSmall1] = m.Registers[m.ArgSmall0]
	}
}
