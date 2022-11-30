# Constants
#Copied from pkg/constants/constants.go
from enum import Enum

class Register(Enum):
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

class Interrupts(Enum):
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
}
