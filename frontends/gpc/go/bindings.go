package main

import "C"
import (
	"log"
	"math"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/vm"
	"unsafe"
)

var gpc vm.VM
var gpcInteruptChannel chan constants.Interrupt = make(chan constants.Interrupt)
var gpcSubbedInterruptChannel chan constants.Interrupt = make(chan constants.Interrupt)

func main() {}

//export Init
func Init(program_bytes *C.char, code_length C.int) {

	vm.InitVM(
		&gpc,
		C.GoBytes(unsafe.Pointer(program_bytes), code_length),
		gpcInteruptChannel,
		gpcSubbedInterruptChannel,
		true,
		false,
	)

	log.Println("VM Created")

}

//export GetInterrupt
func GetInterrupt() C.uint {

	if len(gpc.InterruptArray) > 0 {
		x := gpc.InterruptArray[len(gpc.InterruptArray)-1]
		gpc.InterruptArray = gpc.InterruptArray[:len(gpc.InterruptArray)-1]

		return C.uint(x)
	} else {
		return C.uint(math.MaxUint32)
	}

}

//export SendInterrupt
func SendInterrupt(i C.uint) {

	if gpc.Subscribed(constants.Interrupt(i)) {

		gpc.SubbedInterruptArray = append(gpc.SubbedInterruptArray, constants.Interrupt(i))

	}

}

//export GetRegister
func GetRegister(r C.uint) C.uint {

	x := C.uint(gpc.Registers[r])

	return x

}

//export GetBuffer
func GetBuffer(b C.uint) *C.char {

	//Convert to C.char array

	char_array := []rune{}

	for i := 0; i < 128; i++ {

		if constants.Register(b) == constants.RVideoText {
			char_array = append(char_array, rune(gpc.TextBuffer[i]))
		} else {
			char_array = append(char_array, rune(gpc.DataBuffer[i]))
		}

	}

	return C.CString(string(char_array))

}

//export SetRegister
func SetRegister(r C.uint, v C.uint) {

	gpc.Registers[r] = uint32(v)

}

//export IsSubscribed
func IsSubscribed(i C.uint) C.uint {

	x := gpc.Subscribed(constants.Interrupt(i))

	if x {
		return 1
	} else {
		return 0
	}

}

//export IsFinished
func IsFinished() C.uint {

	x := gpc.Finished

	if x {
		return 1
	} else {
		return 0
	}

}

//export Step
func Step() {

	gpc.Step()

}

//export GetCurrentInstruction
func GetCurrentInstruction() C.uint {

	return C.uint(gpc.Opcode)

}

//export GetArgs
func GetArgs() C.uint {

	return C.uint(gpc.ArgLarge)

}
