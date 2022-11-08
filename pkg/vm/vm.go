package vm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	comp "sccreeper/govm/pkg/compiler"
	c "sccreeper/govm/pkg/constants"
	"sccreeper/govm/pkg/util"
	"sync"
	"time"
)

// General purpose VM backend

const (
	_MemSize                uint32 = 65536 // 2 ^ 16
	_SubscribableInterrupts uint16 = 22
	RegisterCount           uint16 = 52
)

type VM struct {
	MemArray   [_MemSize]byte
	Registers  [RegisterCount + 1]uint32 //float32 or uint32
	DataBuffer [128]byte
	TextBuffer [128]byte

	InterruptTable [_SubscribableInterrupts]uint32

	CurrentInstruction []byte
	Opcode             c.Instruction
	PrevNull           bool
	ProgramBounds      uint32
	Finished           bool

	ArgSmall0 uint16
	ArgSmall1 uint16
	ArgLarge  uint32

	InterruptChannel chan c.Interrupt
	RegisterSync     sync.Mutex
}

// Initialise VM and registers, load code into "memory" etc.
func InitVM(machine *VM, vm_program []byte, interrupt_channel chan c.Interrupt) error {

	if len(vm_program) > int(_MemSize) {
		return errors.New("program too large")
	}

	//Extract program start index

	program_start_index := binary.LittleEndian.Uint32(vm_program[12:])
	interrupt_start_index := binary.LittleEndian.Uint32(vm_program[8:12])

	//Init vars + registers
	machine.Registers[c.RProgramCounter] = program_start_index + comp.StackSize
	machine.CurrentInstruction = vm_program[program_start_index : program_start_index+comp.InstructionLength]
	machine.InterruptChannel = interrupt_channel
	machine.Finished = false
	machine.ProgramBounds = comp.StackSize + uint32(len(vm_program[:len(vm_program)-int(comp.PadSize)]))

	//Interrupt table

	for _, v := range util.SliceChunks(vm_program[interrupt_start_index:program_start_index-comp.PadSize], 6) {

		//log.Println(current_bytes)

		interrupt := c.Interrupt(binary.LittleEndian.Uint16(v[:2]))
		jump_block_addr := binary.LittleEndian.Uint32(v[2:])

		machine.InterruptTable[interrupt] = jump_block_addr

	}

	//Copy program into memory
	copy(machine.MemArray[comp.StackSize:], vm_program[:len(vm_program)-int(comp.PadSize)])

	return nil

}

func (m *VM) Run() {

	for {

		m.RegisterSync.Lock()

		//Interrupts
		select {
		case x := <-m.InterruptChannel:
			m.subbed_interrupt(x)
		default:

		}

		//m.RegisterSync.Lock()

		m.CurrentInstruction = m.MemArray[m.Registers[c.RProgramCounter] : m.Registers[c.RProgramCounter]+comp.InstructionLength]
		m.Opcode = c.Instruction(m.CurrentInstruction[0])

		m.ArgSmall0 = binary.LittleEndian.Uint16(m.CurrentInstruction[1:3])
		m.ArgSmall1 = binary.LittleEndian.Uint16(m.CurrentInstruction[3:5])
		m.ArgLarge = binary.LittleEndian.Uint32(m.CurrentInstruction[1:5])

		//If it is null itn, could be end of program or end of call block
		if m.Opcode == 0 && m.ArgLarge == 0 {

			if m.Registers[c.RProgramCounter] >= m.ProgramBounds {
				m.Finished = true
				break
			}

			fmt.Println("null itn")

			//Set program pointer to call stack pointer
			m.Registers[c.RProgramCounter] = binary.BigEndian.Uint32(m.MemArray[m.Registers[c.RCallStackPointer] : m.Registers[c.RCallStackPointer]+4])

			//Clear previous pointer from call stack
			copy(m.MemArray[m.Registers[c.RStackPointer]:m.Registers[c.RStackPointer]+4], []byte{0, 0, 0, 0})

			m.Registers[c.RCallStackPointer]--
			m.RegisterSync.Unlock()

			continue

		}

		switch m.Opcode {
		//Handle push and pop instructions
		case c.IPush:
			m.push_stack()
		case c.IPop:
			m.pop_stack()

		case c.IMove:
			m.move()

		case c.ICall:
			m.call()
			m.RegisterSync.Unlock()

		case c.IConditionalCall:
			m.conditional_call()
			m.RegisterSync.Unlock()
			continue

		case c.IJump:
			m.jump()
			m.RegisterSync.Unlock()
			continue

		case c.IConditionalJump:
			m.conditional_jump()
			m.RegisterSync.Unlock()
			continue

			// Load & store
		case c.ILoad:
			m.load()

		case c.IStore:
			m.store()

		//Math
		case c.IAdd:
			m.Registers[c.RAccumulator] = m.Registers[m.ArgSmall0] + m.Registers[m.ArgSmall1]
		case c.IMultiply:
			m.Registers[c.RAccumulator] = m.Registers[m.ArgSmall0] * m.Registers[m.ArgSmall1]
		case c.ISubtract:
			m.Registers[c.RAccumulator] = m.Registers[m.ArgSmall0] - m.Registers[m.ArgSmall1]
		case c.IDivide:
			m.Registers[c.RAccumulator] = uint32(m.Registers[m.ArgSmall0] / m.Registers[m.ArgSmall1])
		case c.ISquareRoot:
			m.Registers[c.RAccumulator] = uint32(math.Sqrt(float64(m.Registers[m.ArgSmall0])))
		case c.IIncrement:
			m.Registers[m.ArgSmall0]++
		case c.IDecrement:
			m.Registers[m.ArgSmall0]--
		case c.IInvert:
			m.Registers[m.ArgSmall0] = ^m.Registers[m.ArgSmall0]

		//Logic

		case c.IGreaterThan:
			if m.Registers[m.ArgSmall0] > m.Registers[m.ArgSmall1] {
				m.Registers[c.RAccumulator] = math.MaxUint32
			} else {
				m.Registers[c.RAccumulator] = 0
			}
		case c.ILessThan:
			if m.Registers[m.ArgSmall0] > m.Registers[m.ArgSmall1] {
				m.Registers[c.RAccumulator] = math.MaxUint32
			} else {
				m.Registers[c.RAccumulator] = 0
			}

		case c.IEquals:
			if m.Registers[m.ArgSmall0] == m.Registers[m.ArgSmall1] {
				m.Registers[c.RAccumulator] = math.MaxUint32
			} else {
				m.Registers[c.RAccumulator] = 0
			}
		case c.INotEquals:
			if m.Registers[m.ArgSmall0] != m.Registers[m.ArgSmall1] {
				m.Registers[c.RAccumulator] = math.MaxUint32
			} else {
				m.Registers[c.RAccumulator] = 0
			}

		//Bitwise operators

		case c.IAnd:
			m.Registers[c.RAccumulator] = m.Registers[m.ArgSmall0] & m.Registers[m.ArgSmall1]
		case c.IOr:
			m.Registers[c.RAccumulator] = m.Registers[m.ArgSmall0] | m.Registers[m.ArgSmall1]
		case c.IXor:
			m.Registers[c.RAccumulator] = m.Registers[m.ArgSmall0] ^ m.Registers[m.ArgSmall1]

		case c.IShiftLeft:
			m.Registers[c.RAccumulator] = m.Registers[m.ArgSmall0] << m.Registers[m.ArgSmall1]
		case c.IShiftRight:
			m.Registers[c.RAccumulator] = m.Registers[m.ArgSmall0] >> m.Registers[m.ArgSmall1]

		//Other
		case c.ICallInterrupt:
			m.called_interrupt()
		case c.IHalt:
			time.Sleep(time.Duration(m.Registers[m.ArgSmall0]) * time.Millisecond)

		}

		m.Registers[c.RProgramCounter] += comp.InstructionLength

		m.RegisterSync.Unlock()

	}

	m.RegisterSync.Unlock()
}
