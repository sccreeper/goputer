// Package for defining interupt, register, and instruction constants.

package constants

func init() {

	InterruptIntsReversed = make(map[Interrupt]string)

	for k, v := range InterruptInts {
		InterruptIntsReversed[v] = k
	}

}

// Custom types
type Interrupt uint16
type Register uint16

type Instruction uint8
type InstructionFlag uint8

type DefType uint8

type SoundWaveType uint32

var InterruptIntsReversed map[Interrupt]string

//Maps are used by compiler

var InterruptInts = map[string]Interrupt{

	"ss":  0, //Stop sound
	"sf":  1, //Flush sound registers
	"va":  2, //Render area
	"vp":  3, //Render polygon
	"vt":  4, //Flush video text
	"vc":  5, //Clear video
	"vi":  6, //Draw image
	"vl":  7, //Draw a line from vx0,vy0 -> vx1,vy1
	"iof": 8, //Flush IO registers to IO
	"ioc": 9, //Set all IO to 0x0

	//Subscribable interrupts

	"mm":   10, //Mouse move
	"mu":   11, //Mouse up
	"md":   12, //Mouse down
	"io08": 13, //IO on/off 8-15
	"io09": 14,
	"io10": 15,
	"io11": 16,
	"io12": 17,
	"io13": 18,
	"io14": 19,
	"io15": 20,
	"ku":   21, //Key up
	"kd":   22, //Key down
}

var SubscribableInterrupts = map[string]Interrupt{
	"mm":   10, //Mouse move
	"mu":   11, //Mouse up
	"md":   12, //Mouse down
	"io08": 13, //IO on/off 8-15
	"io09": 14,
	"io10": 15,
	"io11": 16,
	"io12": 17,
	"io13": 18,
	"io14": 19,
	"io15": 20,
	"ku":   21, //Key up
	"kd":   22, //Key down
}

// Array with keys in same order as map
var SubscribableInterruptsKeys []string = []string{
	"mm",   //Mouse move
	"mu",   //Mouse up
	"md",   //Mouse down
	"io08", //IO on/off 8-15
	"io09",
	"io10",
	"io11",
	"io12",
	"io13",
	"io14",
	"io15",
	"ku",
	"kd",
}

var InstructionInts = map[string]uint32{

	"mov": 0, //Move
	"jmp": 1, //Jump

	"add": 2, //Basic arethmetic operations
	"mul": 3,
	"div": 4, //Will floor decimal.
	"sub": 5,

	"cndjmp": 6, //Conditional jump.

	"gt": 7, //Greater than and less than
	"lt": 8,

	"or":  9, //Bitwise logic
	"xor": 10,
	"and": 11,

	"inv": 12, //Invert a number bitewise (flip all bits)

	"eq":  13, //Equals and not equals
	"neq": 14,

	"sl": 15, //Shift left and right
	"sr": 16,

	"int": 17, //Syscall interrupt

	"lda": 18, //Load and store from d0 register
	"sta": 19,

	"push": 20, //Push and pop from stack
	"pop":  21,

	"incr": 22, //Increment register and keep value in register
	"decr": 23,

	"hlt": 24, // Halt the CPU for X milliseconds

	"sqrt": 25, //Square root, will floor decimal.

	"call":    26,
	"cndcall": 27,

	"pow": 28,

	"clr": 29,

	"mod": 30, //Mod instruction

	"emi": 31, //Expansion module interact, not an interrupt because it is handled by the core, not frontends.

	"ret":  32, // Return for normal call
	"iret": 33, // Return for interrupt call
}

var RegisterInts = map[string]uint32{

	"r00": 0, //General purpose registers
	"r01": 1,
	"r02": 2,
	"r03": 3,
	"r04": 4,
	"r05": 5,
	"r06": 6,
	"r07": 7,
	"r08": 8,
	"r09": 9,
	"r10": 10,
	"r11": 11,
	"r12": 12,
	"r13": 13,
	"r14": 14,
	"r15": 15,

	"vx0": 16, //Video X and Y registers
	"vy0": 17,
	"vx1": 18,
	"vy1": 19,

	"vc": 20, //Video colour
	"vb": 21, //Video brightness
	"vt": 22, //Video text (Special register, technically a buffer)

	"kc": 23, //Current key being pressed
	"kp": 24, //Is a key being pressed?

	"mx": 25, //Mouse x and y
	"my": 26,
	"mb": 27, //Current mouse button being pressed.

	"st": 28, //Sound tone
	"sv": 29, //Volume

	"a0": 30, //Accumulator
	"d0": 31, //Data register (returns from interrupts and lda sta)

	"stk": 32, //Current stack pointer
	"stz": 33, //Stack "zero" point in memory

	"io00": 34, //IO registers
	"io01": 35,
	"io02": 36,
	"io03": 37,
	"io04": 38,
	"io05": 39,
	"io06": 40,
	"io07": 41,
	"io08": 42,
	"io09": 43,
	"io10": 44,
	"io11": 45,
	"io12": 46,
	"io13": 47,
	"io14": 48,
	"io15": 49,

	"prc": 50, //Program counter /

	"cstk": 51, //Call stack
	"cstz": 52, //Call stack zero

	"dl": 53, //Data length
	"dp": 54,

	"sw": 55, //Sound wave type

}

//Constants for use in runtime

const (
	//Interrupts

	IntSoundStop    Interrupt = 0
	IntSoundFlush   Interrupt = 1
	IntVideoArea    Interrupt = 2
	IntVideoPolygon Interrupt = 3
	IntVideoText    Interrupt = 4
	IntVideoClear   Interrupt = 5
	IntVideoImage   Interrupt = 6
	IntVideoLine    Interrupt = 7
	IntIOFlush      Interrupt = 8
	IntIOClear      Interrupt = 9

	//Subscribable interrupts

	IntMouseMove Interrupt = 10
	IntMouseUp   Interrupt = 11
	IntMouseDown Interrupt = 12
	IntIO08      Interrupt = 13
	IntIO09      Interrupt = 14
	IntIO10      Interrupt = 15
	IntIO11      Interrupt = 16
	IntIO12      Interrupt = 17
	IntIO13      Interrupt = 18
	IntIO14      Interrupt = 19
	IntIO15      Interrupt = 20

	IntKeyboardUp   Interrupt = 21
	IntKeyboardDown Interrupt = 22
)

const (

	//Instructions

	IMove Instruction = 0
	IJump Instruction = 1

	IAdd      Instruction = 2
	IMultiply Instruction = 3
	IDivide   Instruction = 4
	ISubtract Instruction = 5

	IConditionalJump Instruction = 6

	IGreaterThan Instruction = 7
	ILessThan    Instruction = 8

	IOr  Instruction = 9
	IXor Instruction = 10
	IAnd Instruction = 11

	IInvert Instruction = 12

	IEquals    Instruction = 13
	INotEquals Instruction = 14

	IShiftLeft  Instruction = 15
	IShiftRight Instruction = 16

	ICallInterrupt Instruction = 17

	ILoad  Instruction = 18
	IStore Instruction = 19

	IPush Instruction = 20
	IPop  Instruction = 21

	IIncrement Instruction = 22
	IDecrement Instruction = 23

	IHalt Instruction = 24

	ISquareRoot Instruction = 25

	ICall            Instruction = 26
	IConditionalCall Instruction = 27

	IPower Instruction = 28

	IClear Instruction = 29

	IModulo Instruction = 30

	IExpansionModuleInteract Instruction = 31

	ICallReturn          Instruction = 32
	IInterruptCallReturn Instruction = 33
)

var InstructionArgumentCounts map[Instruction][]int = map[Instruction][]int{

	IJump:                    {1},
	IConditionalJump:         {1},
	IInvert:                  {1},
	ICallInterrupt:           {1},
	ILoad:                    {1, 2},
	IStore:                   {1, 2},
	IIncrement:               {1},
	IDecrement:               {1},
	IHalt:                    {1},
	ISquareRoot:              {1},
	ICall:                    {1},
	IConditionalCall:         {1},
	IClear:                   {1},
	IPush:                    {1},
	IPop:                     {1},
	IExpansionModuleInteract: {1},

	IMove:        {2},
	IAdd:         {2},
	IMultiply:    {2},
	IDivide:      {2},
	ISubtract:    {2},
	IGreaterThan: {2},
	ILessThan:    {2},
	IOr:          {2},
	IXor:         {2},
	IAnd:         {2},
	IEquals:      {2},
	INotEquals:   {2},
	IShiftLeft:   {2},
	IShiftRight:  {2},
	IPower:       {2},
	IModulo:      {2},

	ICallReturn:          {0},
	IInterruptCallReturn: {0},
}

// Determines which arguments in an instruction can have immediate values, if any.
// Even if both values are true, in practice every instruction can only have one immediate value, due to the way immediate values are encoded.
// Both values are true in cases where the order of operations matters.
var InstructionImmediates map[Instruction][][]bool = map[Instruction][][]bool{
	IJump:                    {{true}},
	IConditionalJump:         {{true}},
	IInvert:                  {{false}},
	ICallInterrupt:           {{false}},
	ILoad:                    {{false}, {true, true}},
	IStore:                   {{false}, {true, true}},
	IIncrement:               {{false}},
	IDecrement:               {{false}},
	IHalt:                    {{true}},
	ISquareRoot:              {{true}},
	ICall:                    {{true}},
	IConditionalCall:         {{true}},
	IClear:                   {{true}},
	IPush:                    {{true}},
	IPop:                     {{false}},
	IExpansionModuleInteract: {{true}},

	IMove:        {{true, false}},
	IAdd:         {{false, true}},
	IMultiply:    {{false, true}},
	IDivide:      {{true, true}},
	ISubtract:    {{false, true}},
	IGreaterThan: {{false, true}},
	ILessThan:    {{false, true}},
	IOr:          {{false, true}},
	IXor:         {{false, true}},
	IAnd:         {{false, true}},
	IEquals:      {{false, true}},
	INotEquals:   {{false, true}},
	IShiftLeft:   {{false, true}},
	IShiftRight:  {{false, true}},
	IPower:       {{true, true}},
	IModulo:      {{true, true}},

	ICallReturn:          {{false}},
	IInterruptCallReturn: {{false}},
}

const (
	RGeneralPurpose00 Register = 0
	RGeneralPurpose01 Register = 1
	RGeneralPurpose02 Register = 2
	RGeneralPurpose03 Register = 3
	RGeneralPurpose04 Register = 4
	RGeneralPurpose05 Register = 5
	RGeneralPurpose06 Register = 6
	RGeneralPurpose07 Register = 7
	RGeneralPurpose08 Register = 8
	RGeneralPurpose09 Register = 9
	RGeneralPurpose10 Register = 10
	RGeneralPurpose11 Register = 11
	RGeneralPurpose12 Register = 12
	RGeneralPurpose13 Register = 13
	RGeneralPurpose14 Register = 14
	RGeneralPurpose15 Register = 15

	RVideoX0 Register = 16
	RVideoY0 Register = 17
	RVideoX1 Register = 18
	RVideoY1 Register = 19

	RVideoColour     Register = 20
	RVideoBrightness Register = 21
	RVideoText       Register = 22

	RKeyboardCurrent Register = 23
	RKeyboardPressed Register = 24

	RMouseX      Register = 25
	RMouseY      Register = 26
	RMouseButton Register = 27

	RSoundTone   Register = 28
	RSoundVolume Register = 29

	RAccumulator Register = 30
	RData        Register = 31

	RStackPointer     Register = 32
	RStackZeroPointer Register = 33

	RIO00 Register = 34
	RIO01 Register = 35
	RIO02 Register = 36
	RIO03 Register = 37
	RIO04 Register = 38
	RIO05 Register = 39
	RIO06 Register = 40
	RIO07 Register = 41
	RIO08 Register = 42
	RIO09 Register = 43
	RIO10 Register = 44
	RIO11 Register = 45
	RIO12 Register = 46
	RIO13 Register = 47
	RIO14 Register = 48
	RIO15 Register = 49

	RProgramCounter Register = 50

	RCallStackPointer     Register = 51
	RCallStackZeroPointer Register = 52

	RDataLength  Register = 53
	RDataPointer Register = 54

	RSoundWave Register = 55
)

const (
	StringType DefType = 0
	FloatType  DefType = 1
	IntType    DefType = 2
	UintType   DefType = 3
	BytesType  DefType = 4
)

const (
	SWSquare SoundWaveType = 0
	SWSine   SoundWaveType = 1
)

const (
	ItnFlagFirstArgImmediate InstructionFlag = 0b100_00000
	ItnFlagSecondArgImmediate InstructionFlag = 0b010_00000

	InstructionMask byte = 0b00_111111
	FlagMask byte = ^InstructionMask
	InstructionArgImmediateMask uint32 = 0b000000_11_11111111_11111111_11111111
	InstructionArgRegisterMask uint32 = ^InstructionArgImmediateMask
)
