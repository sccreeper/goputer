// Package for defining interupt, register, and instruction constants.

package tables

var interupt_ints = map[string]int32{

	"ss":  0, //Stop sound
	"va":  1, //Render area
	"vp":  2, //Render pixel
	"io":  3, //Flush IO registers to IO
	"ioc": 4, //Set all IO to 0x0

	// Subscribable interrupts

	"mm":   5, //Mouse move
	"mb":   6, //Mouse button
	"io08": 7, //IO on/off 8-15
	"io09": 8,
	"io10": 9,
	"io11": 10,
	"io12": 11,
	"io13": 12,
	"io14": 13,
	"io15": 14,
	"ku":   15, //Key up
	"kd":   16, //Key down
}

var instructions = map[string]int32{

	"mov": 0, //Move
	"jmp": 1, //Jump

	"add": 2, //Basic arethmetic operations
	"mul": 3,
	"div": 4,
	"sub": 5,

	"cndjmp": 6, //Conditional jump.

	"gt": 7, //Greater than and less than
	"lt": 8,

	"or":  10, //Bitwise logic
	"xor": 11,
	"and": 12,

	"inv": 13, //Invert a number bitewise (flip all bits)

	"eq":  14, //Equals and not equals
	"neq": 15,

	"sl": 16, //Shift left and right
	"sr": 17,

	"int": 18, //Syscall interrupt

	"lda": 19, //Load and store from d0 register
	"sta": 20,
}

var registers = map[string]int32{

	"r00": 0,
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
}
