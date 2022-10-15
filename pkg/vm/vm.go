package vm

import "errors"

// General purpose VM backend

const (
	MemSize                uint32 = 16777216 // 2 ^ 24
	SubscribableInterrupts uint32 = 19
	RegisterCount          uint32 = 50
)

type VM struct {
	MemArray  [MemSize]byte
	Registers [RegisterCount]uint32

	InterruptTable [SubscribableInterrupts]uint32
}

// Initialise VM and registers, load code into "memory" etc.
func InitVM(vm_struct *VM, vm_program []byte) error {

	if len(vm_program) > int(MemSize) {
		return errors.New("program too large")
	}

	//Don't need to init registers, everything

	copy(vm_struct.MemArray[:], vm_program)

	return nil

}
