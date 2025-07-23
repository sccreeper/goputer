package main

import "C"
import (
	"log"
	"math"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/vm"
	"unsafe"
)

var py32 vm.VM
var py32InteruptChannel chan constants.Interrupt = make(chan constants.Interrupt)
var py32SubbedInterruptChannel chan constants.Interrupt = make(chan constants.Interrupt)

func main() {}

//export Init
func Init(programBytes *C.char, codeLength C.int) {

	vm.InitVM(
		&py32,
		C.GoBytes(unsafe.Pointer(programBytes), codeLength),
		false,
	)

	log.Println("VM Created")

}

//export GetInterrupt
func GetInterrupt() C.uint {

	if len(py32.InterruptQueue) > 0 {
		var x constants.Interrupt
		x, py32.InterruptQueue = py32.InterruptQueue[0], py32.InterruptQueue[1:]

		return C.uint(x)
	} else {
		return C.uint(math.MaxUint32)
	}

}

//export SendInterrupt
func SendInterrupt(i C.uint) {

	if py32.Subscribed(constants.Interrupt(i)) {

		py32.SubbedInterruptQueue = append(py32.SubbedInterruptQueue, constants.Interrupt(i))

	}

}

//export GetRegister
func GetRegister(r C.uint) C.uint {

	x := C.uint(py32.Registers[r])

	return x

}

//export GetBuffer
func GetBuffer(b C.uint) *C.char {

	//Convert to C.char array

	charArray := []rune{}

	for i := 0; i < 128; i++ {

		if constants.Register(b) == constants.RVideoText {
			charArray = append(charArray, rune(py32.TextBuffer[i]))
		} else {
			charArray = append(charArray, rune(py32.DataBuffer[i]))
		}

	}

	return C.CString(string(charArray))

}

//export SetRegister
func SetRegister(r C.uint, v C.uint) {

	py32.Registers[r] = uint32(v)

}

//export IsSubscribed
func IsSubscribed(i C.uint) C.uint {

	x := py32.Subscribed(constants.Interrupt(i))

	if x {
		return 1
	} else {
		return 0
	}

}

//export IsFinished
func IsFinished() C.uint {

	x := py32.Finished

	if x {
		return 1
	} else {
		return 0
	}

}

//export Cycle
func Cycle() {

	py32.Cycle()

}

//export GetCurrentInstruction
func GetCurrentInstruction() C.uint {

	return C.uint(py32.Opcode)

}

//export GetArgs
func GetArgs() C.uint {

	return C.uint(py32.LongArg)

}
