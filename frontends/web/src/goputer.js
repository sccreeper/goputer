import * as Comlink from "comlink";
import { ShowError } from "./error";

// Wrapper for worker, also contains types that can be cached, the interrupt and register maps.

/** @type {import("./goputer.worker").Goputer} */
const goputer = Comlink.wrap(new Worker(
  new URL("goputer.worker.js", import.meta.url),
  {type: "module"}
));

goputer.onerror = Comlink.proxy(ShowError)

const interruptInts = {
    "ss":  0,  //Stop sound
	"sf":  1,  //Flush sound registers
	"va":  2,  //Render area
	"vp":  3,  //Render polygon
	"vt":  4,  //Flush video text
	"vc":  5,  //Clear video
	"vi":  6,  //Draw image
	"vl":  7,  //Draw a line from vx0,vy0 -> vx1,vy1
	"iof": 8,  //Flush IO registers to IO
	"ioc": 9,  //Set all IO to 0x0
	"vf":  10, // Video flush

	//Subscribable interrupts

	"mm":   11, //Mouse move
	"mu":   12, //Mouse up
	"md":   13, //Mouse down
	"io08": 14, //IO on/off 8-15
	"io09": 15,
	"io10": 16,
	"io11": 17,
	"io12": 18,
	"io13": 19,
	"io14": 20,
	"io15": 21,
	"ku":   22, //Key up
	"kd":   23, //Key down
	"err":  24,
}

const registerInts = {
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

	"ctrl": 56, // Control register
}

const instructionInts = {
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

	"rand": 34,

	"lteq": 35,
	"gteq": 36,

	"pri": 37, // Prevent interrupts
	"eni": 38, // Enable interrupts
}

const instructionArray = {}

for (const itn in instructionInts) {
    instructionArray[instructionInts[itn]] = itn
}

const interruptArray = {}

for (const ipt in interruptInts) {
    interruptArray[interruptInts[ipt]] = ipt
}

export {goputer, registerInts, interruptInts, instructionInts, instructionArray, interruptArray}