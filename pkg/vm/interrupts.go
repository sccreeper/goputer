package vm

import (
	c "sccreeper/goputer/pkg/constants"
)

type InterruptInfo struct {
	Type uint32
}

func (m *VM) subbedInterrupt(i c.Interrupt) {

	m.ArgLarge = m.InterruptTable[i]

}

func (m *VM) Subscribed(i c.Interrupt) bool {

	return m.InterruptTable[i] != 0

}

func (m *VM) calledInterrupt() {

	if c.Interrupt(m.ArgSmall0) == c.IntIOClear {

		//Set all IO registers to zero

		for i := c.RIO08; i == c.RIO15; i++ {
			m.Registers[i] = 0
		}

	}

	m.InterruptQueue = append(m.InterruptQueue, c.Interrupt(m.ArgSmall0))

}
