package vm

import (
	"errors"
)

type VMHook int

const (
	HookCycle  VMHook = 0
	HookFinish VMHook = 1
)

var HookNames []string = []string{
	"cycle",
	"finish",
}

const hookCount = 2

func (m *VM) AddHook(name string, event VMHook, listener func()) error {

	if _, exists := m.Hooks[event][name]; exists {
		return errors.New("hook already exists")
	}

	m.Hooks[event][name] = listener

	return nil

}

func (m *VM) RemoveHook(name string, event VMHook) {
	delete(m.Hooks[event], name)
}

func (m *VM) CallHooks(event VMHook) {

	for _, v := range m.Hooks[event] {
		v()
	}

}
