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
var py32StepChannel chan bool = make(chan bool)

func main() {}

//export Init
func Init(program_bytes *C.char, code_length C.int) {

	vm.InitVM(
		&py32,
		C.GoBytes(unsafe.Pointer(program_bytes), code_length),
		py32InteruptChannel,
		py32SubbedInterruptChannel,
		true,
		py32StepChannel,
	)

	log.Println("VM Created")

}

//export Run
func Run() {
	go py32.Run()
}

//export GetInterrupt
func GetInterrupt() C.uint {

	select {
	case x := <-py32InteruptChannel:
		return C.uint(x)
	default:
		return C.uint(math.MaxUint32)
	}

}

//export SendInterrupt
func SendInterrupt(i C.uint) {

	if py32.Subscribed(constants.Interrupt(i)) {
		py32SubbedInterruptChannel <- constants.Interrupt(i)
	}

}

//export GetRegister
func GetRegister(r C.uint) C.uint {

	py32.RegisterSync.Lock()
	x := C.uint(py32.Registers[r])
	py32.RegisterSync.Unlock()

	return x

}

//export GetBuffer
func GetBuffer(b C.uint) *C.char {

	//Convert to C.char array

	char_array := []rune{}

	for i := 0; i < 128; i++ {

		if constants.Register(b) == constants.RVideoText {
			char_array = append(char_array, rune(py32.TextBuffer[i]))
		} else {
			char_array = append(char_array, rune(py32.DataBuffer[i]))
		}

	}

	return C.CString(string(char_array))

}

//export SetRegister
func SetRegister(r C.uint, v C.uint) {

	py32.RegisterSync.Lock()
	py32.Registers[r] = uint32(v)
	py32.RegisterSync.Unlock()

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

//export Step
func Step() {

	py32StepChannel <- true

}

//export GetCurrentInstruction
func GetCurrentInstruction() C.uint {

	return C.uint(py32.Opcode)

}

//export GetArgs
func GetArgs() C.uint {

	return C.uint(py32.ArgLarge)

}
