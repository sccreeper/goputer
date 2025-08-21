package vm

import (
	"encoding/binary"
	c "sccreeper/goputer/pkg/constants"
)

func (m *VM) move() {

	var moveLength byte

	if m.Registers[c.RDataLength] > 4 || m.Registers[c.RDataLength] == 0 {
		moveLength = 4
	} else {
		moveLength = byte(m.Registers[c.RDataLength])
	}

	if m.IsImmediate {

		if m.RightArg == uint16(c.RVideoText) || m.RightArg == uint16(c.RData) {

			var val [4]byte

			binary.LittleEndian.PutUint32(val[:], m.LeftArgVal)

			if m.RightArg == uint16(c.RVideoText) {

				for i := 0; i < int(moveLength); i++ {
					m.TextBuffer[i] = val[i]
				}

			} else if m.RightArg == uint16(c.RData) {

				for i := 0; i < int(moveLength); i++ {
					m.DataBuffer[i] = val[i]
				}
			}

		} else {
			m.Registers[m.RightArg] = m.LeftArgVal
		}

	} else {

		if m.LeftArg == uint16(c.RData) && m.RightArg == uint16(c.RVideoText) {

			copy(m.TextBuffer[:m.Registers[c.RDataLength]], m.DataBuffer[:m.Registers[c.RDataLength]])

		} else if m.LeftArg == uint16(c.RVideoText) && m.RightArg == uint16(c.RData) {

			copy(m.DataBuffer[:m.Registers[c.RDataLength]], m.TextBuffer[:m.Registers[c.RDataLength]])

		} else if m.LeftArg == uint16(c.RData) || m.LeftArg == uint16(c.RVideoText) { // buffer -> register

			if m.LeftArg == uint16(c.RData) {
				m.Registers[m.RightArg] = binary.LittleEndian.Uint32(m.DataBuffer[:4])
			} else {
				m.Registers[m.RightArg] = binary.LittleEndian.Uint32(m.TextBuffer[:4])
			}

		} else if m.RightArg == uint16(c.RData) || m.RightArg == uint16(c.RVideoText) { // register -> buffer

			var val [4]byte
			binary.LittleEndian.PutUint32(val[:], m.Registers[m.LeftArg])

			if m.RightArg == uint16(c.RData) {
				for i := 0; i < int(moveLength); i++ {
					m.DataBuffer[i] = val[i]
				}
			} else {
				for i := 0; i < int(moveLength); i++ {
					m.TextBuffer[i] = val[i]
				}
			}

		} else if m.LeftArg < MaxRegister && m.RightArg < MaxRegister {
			m.Registers[m.RightArg] = m.Registers[m.LeftArg]
		}
	}
}
