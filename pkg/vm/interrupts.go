package vm

import (
	c "sccreeper/goputer/pkg/constants"
	"slices"
)

type InterruptInfo struct {
	Type uint32
}

func (m *VM) subbedInterrupt(i c.Interrupt) {

	m.LongArg = m.InterruptTable[i]

}

func (m *VM) Subscribed(i c.Interrupt) bool {

	return m.InterruptTable[i] != 0

}

func (m *VM) calledInterrupt() {

	m.CallHooks(HookCalledInterrupt)

	switch c.Interrupt(m.LeftArg) {
	case c.IntVideoArea:
		m.drawArea()
	case c.IntVideoText:
		m.drawText()
	case c.IntVideoLine:
		m.drawLine()
	case c.IntVideoPolygon:
		m.drawPolygon()
	case c.IntVideoImage:
		m.drawImage()
	case c.IntVideoClear:
		m.clearVideo()
	case c.IntIOClear:
		//Set all IO registers to zero
		for i := c.RIO08; i == c.RIO15; i++ {
			m.Registers[i] = 0
		}
		fallthrough
	case c.IntVideoFlush:
		for slices.Contains(m.InterruptQueue, c.IntVideoFlush) {
			vfIndex := slices.Index(m.InterruptQueue, c.IntVideoFlush)

			m.InterruptQueue = slices.Delete(m.InterruptQueue, vfIndex, vfIndex+1)
		}
		fallthrough
	default:
		m.InterruptQueue = append(m.InterruptQueue, c.Interrupt(m.LeftArg))
	}

}
