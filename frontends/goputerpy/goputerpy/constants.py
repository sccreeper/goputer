# Constants
#Copied from goputer/pkg/constants/constants.go
from enum import IntEnum

class Register(IntEnum):
    RGeneralPurpose00  = 0
    RGeneralPurpose01  = 1
    RGeneralPurpose02  = 2
    RGeneralPurpose03  = 3
    RGeneralPurpose04  = 4
    RGeneralPurpose05  = 5
    RGeneralPurpose06  = 6
    RGeneralPurpose07  = 7
    RGeneralPurpose08  = 8
    RGeneralPurpose09  = 9
    RGeneralPurpose10  = 10
    RGeneralPurpose11  = 11
    RGeneralPurpose12  = 12
    RGeneralPurpose13  = 13
    RGeneralPurpose14  = 14
    RGeneralPurpose15  = 15

    RVideoX0  = 16
    RVideoY0  = 17
    RVideoX1  = 18
    RVideoY1  = 19

    RVideoColour      = 20
    RVideoBrightness  = 21
    RVideoText        = 22

    RKeyboardCurrent  = 23
    RKeyboardPressed  = 24

    RMouseX       = 25
    RMouseY       = 26
    RMouseButton  = 27

    RSoundTone    = 28
    RSoundVolume  = 29

    RAccumulator  = 30
    RData         = 31

    RStackPointer      = 32
    RStackZeroPointer  = 33

    RIO00  = 34
    RIO01  = 35
    RIO02  = 36
    RIO03  = 37
    RIO04  = 38
    RIO05  = 39
    RIO06  = 40
    RIO07  = 41
    RIO08  = 42
    RIO09  = 43
    RIO10  = 44
    RIO11  = 45
    RIO12  = 46
    RIO13  = 47
    RIO14  = 48
    RIO15  = 49

    RProgramCounter  = 50

    RCallStackPointer      = 51
    RCallStackZeroPointer  = 52

    RDataLength   = 53
    RDataPointer  = 54

    RSoundWave  = 55

class Interrupt(IntEnum):
    IntSoundStop   = 0
    IntSoundFlush  = 1
    IntVideoArea   = 2
    IntVideoPixel  = 3
    IntVideoText   = 4
    IntVideoClear  = 5
    IntVideoLine   = 6
    IntIOFlush     = 7
    IntIOClear     = 8

    #Subscribable s

    IntMouseMove  = 9
    IntMouseUp    = 10
    IntMouseDown  = 11
    IntIO08       = 12
    IntIO09       = 13
    IntIO10       = 14
    IntIO11       = 15
    IntIO12       = 16
    IntIO13       = 17
    IntIO14       = 18
    IntIO15       = 19

    IntKeyboardUp    = 20
    IntKeyboardDown  = 21

class SoundWave(IntEnum):
    SWSquare = 0
    SWSine = 1

Instructions = {

	"mov": 0, #Move
	"jmp": 1, #Jump

	"add": 2, #Basic arethmetic operations
	"mul": 3,
	"div": 4,
	"sub": 5,

	"cndjmp": 6, #Conditional jump.

	"gt": 7, #Greater than and less than
	"lt": 8,

	"or":  9, #Bitwise logic
	"xor": 10,
	"and": 11,

	"inv": 12, #Invert a number bitewise (flip all bits)

	"eq":  13, #Equals and not equals
	"neq": 14,

	"sl": 15, #Shift left and right
	"sr": 16,

	"int": 17, #Syscall interrupt

	"lda": 18, #Load and store from d0 register
	"sta": 19,

	"push": 20, #Push and pop from stack
	"pop":  21,

	"incr": 22, #Increment register and keep value in register
	"decr": 23,

	"hlt": 24, # Halt the CPU for X milliseconds

	"sqrt": 25, #Square root, will round to nearest uint it isn't a float instruction

	"call":    26,
	"cndcall": 27,

	"pow": 28,

	"clr": 29,
     
	"mod" : 30,
    "emi" : 31,
}

SingleArgInstructions = [
	Instructions["jmp"],
	Instructions["cndjmp"],
	Instructions["inv"],
	Instructions["call"],
	Instructions["lda"],
	Instructions["sta"],
	Instructions["incr"],
	Instructions["decr"],
	Instructions["hlt"],
	Instructions["sqrt"],
	Instructions["call"],
	Instructions["cndcall"],
	Instructions["clr"],
    Instructions["emi"],
]

InterruptInts = {

	"ss":  0, #Stop sound
	"sf":  1, #Flush sound registers
	"va":  2, #Render area
	"vp":  3, #Render pixel
	"vt":  4, #Flush video text
	"vc":  5, #Clear video
	"vl":  6, #Draw a line from vx0,vy0 -> vx1,vy1
	"iof": 7, #Flush IO registers to IO
	"ioc": 8, #Set all IO to 0x0

	#Subscribable interrupts

	"mm":   9,  #Mouse move
	"mu":   10, #Mouse up
	"md":   11, #Mouse down
	"io08": 12, #IO on/off 8-15
	"io09": 13,
	"io10": 14,
	"io11": 15,
	"io12": 16,
	"io13": 17,
	"io14": 18,
	"io15": 19,
	"ku":   20, #Key up
	"kd":   21, #Key down
}

RegisterInts = {

	"r00": 0, #General purpose registers
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

	"vx0": 16, #Video X and Y registers
	"vy0": 17,
	"vx1": 18,
	"vy1": 19,

	"vc": 20, #Video colour
	"vb": 21, #Video brightness
	"vt": 22, #Video text (Special register, technically a buffer)

	"kc": 23, #Current key being pressed
	"kp": 24, #Is a key being pressed?

	"mx": 25, #Mouse x and y
	"my": 26,
	"mb": 27, #Current mouse button being pressed.

	"st": 28, #Sound tone
	"sv": 29, #Volume

	"a0": 30, #Accumulator
	"d0": 31, #Data register (returns from interrupts and lda sta)

	"stk": 32, #Current stack pointer
	"stz": 33, #Stack "zero" point in memory

	"io00": 34, #IO registers
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

	"prc": 50, #Program counter /

	"cstk": 51, #Call stack
	"cstz": 52, #Call stack zero

	"dl": 53, #Data length
	"dp": 54,

	"sw": 55, #Sound wave type
}

RegisterStrings = {}

for k in RegisterInts:
    RegisterStrings[RegisterInts[k]] = k

InterruptStrings = {}

for k in InterruptInts:
    InterruptStrings[InterruptInts[k]] = k

InstructionStrings = {}

for k in Instructions:
	InstructionStrings[Instructions[k]] = k

