package vm

import (
	c "sccreeper/goputer/pkg/constants"
)

type InterruptInfo struct {
	Type uint32
}

func (m *VM) subbed_interrupt(i c.Interrupt) {

	m.ArgLarge = m.InterruptTable[i]

}

func (m *VM) Subscribed(i c.Interrupt) bool {

	if m.InterruptTable[i] == 0 {
		return false
	} else {

		return true

	}

}

func (m *VM) called_interrupt() {

	if c.Interrupt(m.ArgSmall0) == c.IntIOClear {

		//Set all IO registers to zero

		for i := c.RIO08; i == c.RIO15; i++ {
			m.Registers[i] = 0
		}

	}

	m.InterruptChannel <- c.Interrupt(m.ArgSmall0)
}
