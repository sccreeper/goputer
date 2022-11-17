// Package for defining interupt, register, and instruction constants.

package constants

// Custom types
type Interrupt uint16
type Register uint16

type Instruction uint8

type DefType uint8

//Maps are used by compiler

var InterruptInts = map[string]Interrupt{

	"ss":  0, //Stop sound
	"sf":  1, //Flush sound registers
	"va":  2, //Render area
	"vp":  3, //Render pixel
	"vt":  4, //Flush video text
	"vc":  5, //Clear video
	"vl":  6, //Draw a line from vx0,vy0 -> vx1,vy1
	"iof": 7, //Flush IO registers to IO
	"ioc": 8, //Set all IO to 0x0

	//Subscribable interrupts

	"mm":   9,  //Mouse move
	"mu":   10, //Mouse up
	"md":   11, //Mouse down
	"io08": 12, //IO on/off 8-15
	"io09": 13,
	"io10": 14,
	"io11": 15,
	"io12": 16,
	"io13": 17,
	"io14": 18,
	"io15": 19,
	"ku":   20, //Key up
	"kd":   21, //Key down
}

var SubscribableInterrupts = map[string]Interrupt{
	"mm":   9,  //Mouse move
	"mu":   10, //Mouse up
	"md":   11, //Mouse down
	"io08": 12, //IO on/off 8-15
	"io09": 13,
	"io10": 14,
	"io11": 15,
	"io12": 16,
	"io13": 17,
	"io14": 18,
	"io15": 19,
	"ku":   20, //Key up
	"kd":   21, //Key down
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
	"div": 4,
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

	"sqrt": 25, //Square root, will round to nearest uint it isn't a float instruction

	"call":    26,
	"cndcall": 27,
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
}

//Constants for use in runtime

const (
	//Interrupts

	IntSoundStop  Interrupt = 0
	IntSoundFlush Interrupt = 1
	IntVideoArea  Interrupt = 2
	IntVideoPixel Interrupt = 3
	IntVideoText  Interrupt = 4
	IntVideoClear Interrupt = 5
	IntVideoLine  Interrupt = 6
	IntIOFlush    Interrupt = 7
	IntIOClear    Interrupt = 8

	//Subscribable interrupts

	IntMouseMove Interrupt = 9
	IntMouseUp   Interrupt = 10
	IntMouseDown Interrupt = 12
	IntIO08      Interrupt = 13
	IntIO09      Interrupt = 14
	IntIO10      Interrupt = 15
	IntIO11      Interrupt = 16
	IntIO12      Interrupt = 17
	IntIO13      Interrupt = 18
	IntIO14      Interrupt = 19
	IntIO15      Interrupt = 20

	IntKeyboardUp   Interrupt = 20
	IntKeyboardDown Interrupt = 21
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
)

// Instructions that take a single 32 bit arg, as opposed to 2x16bit args
var SingleArgInstructions = []Instruction{

	IJump,
	IConditionalJump,
	IInvert,
	ICallInterrupt,
	ILoad,
	IStore,
	IIncrement,
	IDecrement,
	IHalt,
	ISquareRoot,
	ICall,
	IConditionalCall,
	IClear,
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
	RIO08 Register = 32
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
)

const (
	StringType DefType = 0
	FloatType  DefType = 1
	IntType    DefType = 2
	UintType   DefType = 3
	BytesType  DefType = 4
)
