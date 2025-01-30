package vm

import (
	"encoding/binary"
	"errors"
	"log"
	"math"
	comp "sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/constants"
	c "sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/expansions"
	"sccreeper/goputer/pkg/util"
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

	InterruptQueue       []c.Interrupt
	SubbedInterruptQueue []c.Interrupt
	HandlingInterrupt    bool

	ExecutionPaused    bool
	ExecutionPauseTime int64

	ExpansionsSupported bool
}

// Initialize VM and registers, load code into "memory" etc.
func InitVM(machine *VM, vmProgram []byte, expansionsSupported bool) error {

	if len(vmProgram)-4 > int(_MemSize) {
		return errors.New("program too large")
	}

	//Extract program start index

	vmProgram = vmProgram[4:]

	program_start_index := binary.LittleEndian.Uint32(vmProgram[12:])
	interrupt_start_index := binary.LittleEndian.Uint32(vmProgram[8:12])

	//Init vars + registers
	machine.Registers[c.RProgramCounter] = program_start_index + comp.StackSize
	machine.CurrentInstruction = vmProgram[program_start_index : program_start_index+comp.InstructionLength]
	machine.Finished = false
	machine.ProgramBounds = comp.StackSize + uint32(len(vmProgram[:len(vmProgram)-int(comp.PadSize)]))
	machine.Registers[c.RVideoBrightness] = 255
	machine.ExpansionsSupported = expansionsSupported

	machine.Registers[c.RCallStackZeroPointer] = comp.StackSize - comp.CallStackSize
	machine.Registers[c.RCallStackPointer] = machine.Registers[c.RCallStackZeroPointer]

	machine.Registers[c.RStackZeroPointer] = 0
	machine.Registers[c.RStackPointer] = 0

	machine.InterruptQueue = []c.Interrupt{}
	machine.SubbedInterruptQueue = []c.Interrupt{}
	machine.HandlingInterrupt = false

	//Interrupt table

	for _, v := range util.SliceChunks(vmProgram[interrupt_start_index:program_start_index-comp.PadSize], 6) {

		//log.Println(current_bytes)

		interrupt := c.Interrupt(binary.LittleEndian.Uint16(v[:2]))
		jump_block_addr := binary.LittleEndian.Uint32(v[2:])

		machine.InterruptTable[interrupt] = jump_block_addr

	}

	//Copy program into memory
	copy(machine.MemArray[comp.StackSize:], vmProgram[:len(vmProgram)-int(comp.PadSize)])

	// Load expansions

	if expansionsSupported {
		expansions.LoadExpansions()
	} else {
		log.Println("Expansions are disabled for this frontend")
	}

	return nil

}

func (m *VM) Cycle() {

	// Stop if the program has terminated

	if m.Finished {
		return
	}

	// If we are in the middle of a halt, pause then continue

	if m.ExecutionPaused {

		if time.Now().UnixMilli()-m.ExecutionPauseTime >= int64(m.Registers[m.ArgSmall0]) {
			m.ExecutionPaused = false
			m.Registers[c.RProgramCounter] += comp.InstructionLength
			return
		} else {
			return
		}

	}

	// Get arguments from memory

	m.CurrentInstruction = m.MemArray[m.Registers[c.RProgramCounter] : m.Registers[c.RProgramCounter]+comp.InstructionLength]
	m.Opcode = c.Instruction(m.CurrentInstruction[0])

	m.ArgSmall0 = binary.LittleEndian.Uint16(m.CurrentInstruction[1:3])
	m.ArgSmall1 = binary.LittleEndian.Uint16(m.CurrentInstruction[3:5])
	m.ArgLarge = binary.LittleEndian.Uint32(m.CurrentInstruction[1:5])

	//Interrupts

	if !m.HandlingInterrupt && len(m.SubbedInterruptQueue) > 0 {

		// Pop from queue
		var i constants.Interrupt
		i, m.SubbedInterruptQueue = m.SubbedInterruptQueue[0], m.SubbedInterruptQueue[1:]

		// Frontends should do the checking but this is just to be sure.
		if m.Subscribed(i) {
			m.HandlingInterrupt = true

			m.subbedInterrupt(i)
			m.call()
			return
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
			m.HandlingInterrupt = false
			m.popCall()
			return
		}

	}

	switch m.Opcode {
	//Handle push and pop instructions
	case c.IPush:
		m.pushStack()
	case c.IPop:
		m.popStack()

	case c.IMove:
		m.move()

	case c.ICall:
		m.call()
		return

	case c.IConditionalCall:
		if m.conditionalCall() {
			return
		}

	case c.IJump:
		m.jump()
		return

	case c.IConditionalJump:
		if m.conditionalJump() {
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
		m.shiftLeft()
	case c.IShiftRight:
		m.shiftRight()

	//Other
	case c.ICallInterrupt:
		m.calledInterrupt()
	case c.IHalt:

		m.ExecutionPaused = true
		m.ExecutionPauseTime = time.Now().UnixMilli()
		return

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
