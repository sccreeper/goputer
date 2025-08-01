package vm

import (
	"encoding/binary"
	"errors"
	"math"
	comp "sccreeper/goputer/pkg/compiler"
	c "sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/expansions"
	"sccreeper/goputer/pkg/util"
	"time"
)

// General purpose VM backend

const (
	MemSize          uint32 = VideoBufferSize + 65536 // 2 ^ 16
	MaxRegister    uint16 = 55
	InstructionCount uint16 = 34
	InterruptCount   uint16 = 24
)

type VM struct {
	MemArray   [MemSize]byte
	Registers  [MaxRegister + 1]uint32 //float32 or uint32
	DataBuffer [128]byte
	TextBuffer [128]byte

	InterruptTable [InterruptCount]uint32

	CurrentInstruction []byte
	Opcode             c.Instruction
	PrevNull           bool
	ProgramBounds      uint32
	Finished           bool

	LeftArg uint16
	RightArg uint16
	LongArg  uint32

	LeftArgVal uint32
	RightArgVal uint32
	LongArgVal uint32

	IsImmediate bool
	ImmediateArgIndex int

	InterruptQueue       []c.Interrupt
	SubbedInterruptQueue []c.Interrupt
	HandlingInterrupt    bool

	ExecutionPaused    bool
	ExecutionPauseTime int64

	ExpansionsSupported bool
}

// Initialize VM and registers, load code into "memory" etc.
func NewVM(vmProgram []byte, expansionsSupported bool) (*VM, error) {

	machine := &VM{}

	if len(vmProgram)-4 > int(MemSize) {
		return nil, errors.New("program too large")
	}

	//Extract program start index

	vmProgram = vmProgram[4:]

	var interruptStartIndex uint32 = binary.LittleEndian.Uint32(vmProgram[:4])
	var definitionStartIndex uint32 = binary.LittleEndian.Uint32(vmProgram[4:8])
	// instructionStart = 8:12
	var instructionEntryPoint uint32 = binary.LittleEndian.Uint32(vmProgram[12:16])

	//Init vars + registers
	machine.Registers[c.RProgramCounter] = instructionEntryPoint + comp.MemOffset
	machine.CurrentInstruction = vmProgram[instructionEntryPoint : instructionEntryPoint+comp.InstructionLength]
	machine.Finished = false
	machine.ProgramBounds = comp.MemOffset + uint32(len(vmProgram[:len(vmProgram)-int(comp.PadSize)]))
	machine.Registers[c.RVideoBrightness] = 255
	machine.ExpansionsSupported = expansionsSupported

	machine.Registers[c.RCallStackZeroPointer] = comp.MemOffset - comp.CallStackSize
	machine.Registers[c.RCallStackPointer] = machine.Registers[c.RCallStackZeroPointer]

	machine.Registers[c.RStackZeroPointer] = comp.MemOffset - comp.CallStackSize - comp.DataStackSize
	machine.Registers[c.RStackPointer] = machine.Registers[c.RStackZeroPointer]

	machine.InterruptQueue = make([]c.Interrupt, 0, 16)
	machine.SubbedInterruptQueue = make([]c.Interrupt, 0, 16)
	machine.HandlingInterrupt = false

	//Interrupt table

	for _, v := range util.SliceChunks(vmProgram[interruptStartIndex:definitionStartIndex], 6) {

		//log.Println(current_bytes)

		interrupt := c.Interrupt(binary.LittleEndian.Uint16(v[:2]))
		jumpBlockAddr := binary.LittleEndian.Uint32(v[2:])

		machine.InterruptTable[interrupt] = jumpBlockAddr

	}

	//Copy program into memory
	copy(machine.MemArray[comp.MemOffset:], vmProgram[:len(vmProgram)-int(comp.PadSize)])

	// Load expansions

	if expansionsSupported {
		expansions.LoadExpansions()
	}

	return machine, nil

}

func (m *VM) Cycle() {

	// Stop if the program has terminated

	if m.Finished {
		return
	}

	// If we are in the middle of a halt, pause then continue

	if m.ExecutionPaused {

		if time.Now().UnixMilli()-m.ExecutionPauseTime >= int64(m.LeftArgVal) {
			m.ExecutionPaused = false
			m.Registers[c.RProgramCounter] += comp.InstructionLength
			return
		} else {
			return
		}

	}

	m.IsImmediate = false

	// Get arguments from memory

	m.CurrentInstruction = m.MemArray[m.Registers[c.RProgramCounter] : m.Registers[c.RProgramCounter]+comp.InstructionLength]
	m.Opcode = c.Instruction(m.CurrentInstruction[0] & c.InstructionMask)

	m.LeftArg = binary.LittleEndian.Uint16(m.CurrentInstruction[1:3])
	m.RightArg = binary.LittleEndian.Uint16(m.CurrentInstruction[3:5])
	m.LongArg = binary.LittleEndian.Uint32(m.CurrentInstruction[1:5])

	if (m.CurrentInstruction[0] & byte(c.ItnFlagLongArgImmediate)) == byte(c.ItnFlagLongArgImmediate) {
		
		m.LeftArgVal = uint32(m.LeftArg)
		m.RightArgVal = uint32(m.RightArg)
		m.LongArgVal = m.LongArg

	} else if (m.CurrentInstruction[0] & byte(c.ItnFlagLeftArgImmediate)) != 0 || (m.CurrentInstruction[0] & byte(c.ItnFlagRightArgImmediate)) != 0 {
		m.IsImmediate = true

		immVal := m.LongArg & c.InstructionArgImmediateMask
		immReg := (m.LongArg & c.InstructionArgRegisterMask) >> 26

		if (m.CurrentInstruction[0] & byte(c.ItnFlagLeftArgImmediate)) != 0 {
			m.ImmediateArgIndex = 0
			m.LeftArgVal = immVal
			m.LeftArg = 0
			m.RightArgVal = m.Registers[immReg]
			m.RightArg = uint16(immReg)
		} else {
			m.ImmediateArgIndex = 1
			m.RightArgVal = immVal
			m.RightArg = 0
			m.LeftArgVal = m.Registers[immReg]
			m.LeftArg = uint16(immReg)
		}

		m.LongArgVal = immVal

	} else {
		if m.LeftArg < MaxRegister {
			m.LeftArgVal = m.Registers[m.LeftArg]
		}

		if m.RightArg < MaxRegister {
			m.RightArgVal = m.Registers[m.RightArg]
		}

		m.LongArgVal = m.LongArg
	}

	//Interrupts

	if !m.HandlingInterrupt && len(m.SubbedInterruptQueue) > 0 {

		// Pop from queue
		var i c.Interrupt
		i, m.SubbedInterruptQueue = m.SubbedInterruptQueue[0], m.SubbedInterruptQueue[1:]

		// Frontends should do the checking but this is just to be sure.
		if m.Subscribed(i) {
			m.HandlingInterrupt = true

			m.subbedInterrupt(i)
			m.call()
			return
		}

	}

	// If there is a null instruction, then terminate program.
	// Null instructions should only ever be encountered this way.
	if m.Opcode == 0 && m.LongArg == 0 {

		m.Finished = true
		return

	}

	switch m.Opcode {

	// Handle push and pop instructions

	case c.IPush:
		m.pushStack()
	case c.IPop:
		m.popStack()

	case c.IMove:
		m.move()

	// Control flow

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

	case c.IInterruptCallReturn:
		m.HandlingInterrupt = false
		m.popCall()
		return
	case c.ICallReturn:
		m.popCall()
		return

	// Load & store

	case c.ILoad:
		m.load()

	case c.IStore:
		m.store()

	//Math
	case c.IAdd:
		m.Registers[c.RAccumulator] = m.LeftArgVal + m.RightArgVal
	case c.IMultiply:
		m.Registers[c.RAccumulator] = m.LeftArgVal * m.RightArgVal
	case c.ISubtract:
		m.Registers[c.RAccumulator] = m.LeftArgVal - m.RightArgVal
	case c.IDivide:
		m.Registers[c.RAccumulator] = m.LeftArgVal / m.RightArgVal
	case c.ISquareRoot:
		m.Registers[c.RAccumulator] = uint32(math.Sqrt(float64(m.LeftArgVal)))
	case c.IIncrement:
		m.Registers[m.LeftArg]++
	case c.IDecrement:
		m.Registers[m.LeftArg]--
	case c.IInvert:
		m.Registers[c.RAccumulator] = ^m.LeftArgVal
	case c.IPower:
		if m.LeftArgVal == 10 {
			m.Registers[c.RAccumulator] = uint32(math.Pow10(int(m.RightArgVal)))
		} else {
			m.Registers[c.RAccumulator] = uint32(math.Pow(float64(m.LeftArgVal), float64(m.RightArgVal)))
		}
	case c.IModulo:
		m.Registers[c.RAccumulator] = m.LeftArgVal % m.RightArgVal

	//Logic

	case c.IGreaterThan:
		if m.LeftArgVal > m.RightArgVal {
			m.Registers[c.RAccumulator] = math.MaxUint32
		} else {
			m.Registers[c.RAccumulator] = 0
		}
	case c.ILessThan:
		if m.LeftArgVal < m.RightArgVal {
			m.Registers[c.RAccumulator] = math.MaxUint32
		} else {
			m.Registers[c.RAccumulator] = 0
		}

	case c.IEquals:
		if m.LeftArgVal == m.RightArgVal {
			m.Registers[c.RAccumulator] = math.MaxUint32
		} else {
			m.Registers[c.RAccumulator] = 0
		}
	case c.INotEquals:
		if m.LeftArgVal != m.RightArgVal {
			m.Registers[c.RAccumulator] = math.MaxUint32
		} else {
			m.Registers[c.RAccumulator] = 0
		}

	//Bitwise operators

	case c.IAnd:
		m.Registers[c.RAccumulator] = m.LeftArgVal & m.RightArgVal
	case c.IOr:
		m.Registers[c.RAccumulator] = m.LeftArgVal | m.RightArgVal
	case c.IXor:
		m.Registers[c.RAccumulator] = m.LeftArgVal ^ m.RightArgVal

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
		if m.LongArg != uint32(c.RData) && m.LongArg != uint32(c.RVideoText) {

			m.Registers[m.LongArg] = 0

		} else {

			switch m.LongArg {
			case uint32(c.RData):
				m.DataBuffer = [128]byte{}
			case uint32(c.RVideoText):
				m.TextBuffer = [128]byte{}
			}

		}
	case c.IExpansionModuleInteract:
		if expansions.ModuleExists(m.LongArgVal) && m.ExpansionsSupported {
			data := expansions.Interaction(m.LongArgVal, m.DataBuffer[:])

			m.Registers[c.RDataLength] = uint32(len(data))
			m.Registers[c.RDataPointer] = 0
			copy(m.DataBuffer[:], data)
		}
	}

	m.Registers[c.RProgramCounter] += comp.InstructionLength

}
