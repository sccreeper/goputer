package vm

import (
	"encoding/binary"
	"errors"
	"log"
	"math"
	comp "sccreeper/goputer/pkg/compiler"
	c "sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/expansions"
	"sccreeper/goputer/pkg/util"
	"sync"
	"time"
)

// General purpose VM backend

const (
	_MemSize                uint32 = 65536 // 2 ^ 16
	_SubscribableInterrupts uint16 = 22
	RegisterCount           uint16 = 57
	InstructionCount        uint16 = 32
	InterruptCount          uint16 = 22
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

	InterruptChannel       chan c.Interrupt
	SubbedInterruptChannel chan c.Interrupt
	RegisterSync           sync.Mutex

	InterruptArray       []c.Interrupt
	SubbedInterruptArray []c.Interrupt

	HandlingInterrupt bool
	InterruptQueue    []uint32

	ShouldStep bool

	ExpansionsSupported bool
}

// Initialize VM and registers, load code into "memory" etc.
func InitVM(machine *VM, vm_program []byte, interrupt_channel chan c.Interrupt, subbed_interrupt_channel chan c.Interrupt, should_step bool, expansions_supported bool) error {

	if len(vm_program)-4 > int(_MemSize) {
		return errors.New("program too large")
	}

	//Extract program start index

	vm_program = vm_program[4:]

	program_start_index := binary.LittleEndian.Uint32(vm_program[12:])
	interrupt_start_index := binary.LittleEndian.Uint32(vm_program[8:12])

	//Init vars + registers
	machine.Registers[c.RProgramCounter] = program_start_index + comp.StackSize
	machine.CurrentInstruction = vm_program[program_start_index : program_start_index+comp.InstructionLength]
	machine.InterruptChannel = interrupt_channel
	machine.SubbedInterruptChannel = subbed_interrupt_channel
	machine.Finished = false
	machine.ProgramBounds = comp.StackSize + uint32(len(vm_program[:len(vm_program)-int(comp.PadSize)]))
	machine.Registers[c.RVideoBrightness] = 255
	machine.ExpansionsSupported = expansions_supported

	machine.Registers[c.RCallStackZeroPointer] = comp.StackSize - comp.CallStackSize
	machine.Registers[c.RCallStackPointer] = machine.Registers[c.RCallStackZeroPointer]

	machine.Registers[c.RStackZeroPointer] = 0
	machine.Registers[c.RStackPointer] = 0

	machine.ShouldStep = should_step

	machine.HandlingInterrupt = false
	machine.InterruptQueue = []uint32{}

	machine.InterruptArray = []c.Interrupt{}
	machine.SubbedInterruptArray = []c.Interrupt{}

	//Interrupt table

	for _, v := range util.SliceChunks(vm_program[interrupt_start_index:program_start_index-comp.PadSize], 6) {

		//log.Println(current_bytes)

		interrupt := c.Interrupt(binary.LittleEndian.Uint16(v[:2]))
		jump_block_addr := binary.LittleEndian.Uint32(v[2:])

		machine.InterruptTable[interrupt] = jump_block_addr

	}

	//Copy program into memory
	copy(machine.MemArray[comp.StackSize:], vm_program[:len(vm_program)-int(comp.PadSize)])

	// Load expansions

	if expansions_supported {
		expansions.LoadExpansions()
	} else {
		log.Println("Expansions are disabled for this frontend")
	}

	return nil

}

func (m *VM) Run() {

	if m.ShouldStep {
		panic(errors.New("run called when behaivour is step"))
	} else {
		go func() {

			for {

				if m.Finished {
					break
				}

				m.Cycle()

			}

			close(m.SubbedInterruptChannel)
			close(m.InterruptChannel)

		}()
	}

}

func (m *VM) Step() {
	if m.Finished {
		close(m.SubbedInterruptChannel)
		close(m.InterruptChannel)
		return
	}

	m.Cycle()
}

func (m *VM) Cycle() {

	if m.Finished {
		return
	}

	m.CurrentInstruction = m.MemArray[m.Registers[c.RProgramCounter] : m.Registers[c.RProgramCounter]+comp.InstructionLength]
	m.Opcode = c.Instruction(m.CurrentInstruction[0])

	m.ArgSmall0 = binary.LittleEndian.Uint16(m.CurrentInstruction[1:3])
	m.ArgSmall1 = binary.LittleEndian.Uint16(m.CurrentInstruction[3:5])
	m.ArgLarge = binary.LittleEndian.Uint32(m.CurrentInstruction[1:5])

	//Interrupts

	var i uint32 = math.MaxUint32

	if !m.ShouldStep {
		select {
		case x := <-m.SubbedInterruptChannel:
			i = uint32(x)
		default:

		}

	} else {

		if len(m.SubbedInterruptArray) > 0 {
			i = uint32(m.SubbedInterruptArray[len(m.SubbedInterruptArray)-1])
			m.SubbedInterruptArray = m.SubbedInterruptArray[:len(m.SubbedInterruptArray)-1]
		}
	}

	if i != math.MaxUint32 {
		//Place the interrupts into a "queue"
		if m.HandlingInterrupt {
			if m.Subscribed(c.Interrupt(i)) {
				m.InterruptQueue = append(m.InterruptQueue, uint32(i))
			}
		} else {

			var interrupt c.Interrupt

			if len(m.InterruptQueue) == 0 {
				interrupt = c.Interrupt(i)
			} else {
				interrupt = c.Interrupt(m.InterruptQueue[len(m.InterruptQueue)-1])
				m.InterruptQueue = m.InterruptQueue[:len(m.InterruptQueue)-1]
			}

			if m.Subscribed(interrupt) {
				m.HandlingInterrupt = true
				m.subbed_interrupt(interrupt)
				m.call()
				//fmt.Printf("Interrupt %d\n", x)
				return
			}
		}
	}

	//If it is null itn, could be end of program or end of call block
	if m.Opcode == 0 && m.ArgLarge == 0 {

		//If the next opcode and arg is 0 as well we exit
		//If not we pop from the call stack

		next_instruction := m.MemArray[m.Registers[c.RProgramCounter]+comp.InstructionLength : m.Registers[c.RProgramCounter]+(comp.InstructionLength*2)]

		if next_instruction[0] == 0 && util.AllEqualToX(m.CurrentInstruction[1:5], 0) {
			m.Finished = true
			return
		} else {
			m.pop_call()
			m.HandlingInterrupt = false

			return
		}

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
		return

	case c.IConditionalCall:
		if m.conditional_call() {
			return
		}

	case c.IJump:
		m.jump()
		return

	case c.IConditionalJump:
		if m.conditional_jump() {
			return
		}

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
		m.Registers[c.RAccumulator] = uint32(math.Floor(float64(m.Registers[m.ArgSmall0]) / float64(m.Registers[m.ArgSmall1])))
	case c.ISquareRoot:
		m.Registers[c.RAccumulator] = uint32(math.Sqrt(float64(m.Registers[m.ArgSmall0])))
	case c.IIncrement:
		m.Registers[m.ArgSmall0]++
	case c.IDecrement:
		m.Registers[m.ArgSmall0]--
	case c.IInvert:
		m.Registers[m.ArgSmall0] = ^m.Registers[m.ArgSmall0]
	case c.IPower:
		if m.Registers[m.ArgSmall0] == 10 {
			m.Registers[c.RAccumulator] = uint32(math.Pow10(int(m.Registers[m.ArgSmall1])))
		} else {
			m.Registers[c.RAccumulator] = uint32(math.Pow(float64(m.Registers[m.ArgSmall0]), float64(m.Registers[m.ArgSmall1])))
		}
	case c.IModulo:
		m.Registers[c.RAccumulator] = m.Registers[m.ArgSmall0] % m.Registers[m.ArgSmall1]

	//Logic

	case c.IGreaterThan:
		if m.Registers[m.ArgSmall0] > m.Registers[m.ArgSmall1] {
			m.Registers[c.RAccumulator] = math.MaxUint32
		} else {
			m.Registers[c.RAccumulator] = 0
		}
	case c.ILessThan:
		if m.Registers[m.ArgSmall0] < m.Registers[m.ArgSmall1] {
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
		m.shift_left()
	case c.IShiftRight:
		m.shift_right()

	//Other
	case c.ICallInterrupt:
		m.called_interrupt()
	case c.IHalt:
		time.Sleep(time.Duration(m.Registers[m.ArgSmall0]) * time.Millisecond)
	case c.IClear:
		if m.ArgLarge != uint32(c.RData) && m.ArgLarge != uint32(c.RVideoText) {

			m.Registers[m.ArgLarge] = 0

		} else {

			switch m.ArgLarge {
			case uint32(c.RData):
				m.DataBuffer = [128]byte{}
			case uint32(c.RVideoText):
				m.TextBuffer = [128]byte{}
			}

		}
	case c.IExpansionModuleInteract:
		if expansions.ModuleExists(m.Registers[m.ArgLarge]) && m.ExpansionsSupported {
			data := expansions.Interaction(m.Registers[m.ArgLarge], m.DataBuffer[:])

			m.Registers[c.RDataLength] = uint32(len(data))
			m.Registers[c.RDataPointer] = 0
			copy(m.DataBuffer[:], data)
		}
	}

	m.Registers[c.RProgramCounter] += comp.InstructionLength

}
