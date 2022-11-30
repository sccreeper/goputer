package main

import "C"
import (
	"log"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/vm"
	"unsafe"
)

var py32 vm.VM
var py32InteruptChannel chan constants.Interrupt = make(chan constants.Interrupt)
var py32SubbedInterruptChannel chan constants.Interrupt = make(chan constants.Interrupt)

func main() {}

//export Init
func Init(program_bytes []C.char, code_length C.int) {

	vm.InitVM(
		&py32,
		C.GoBytes(unsafe.Pointer(&program_bytes), code_length),
		py32InteruptChannel,
		py32SubbedInterruptChannel,
	)

	log.Println("VM Created")

}

//export Run
func Run() {
	go py32.Run()
}

//export GetInterrupt
func GetInterrupt() C.ulong {

	select {
	case x := <-py32InteruptChannel:
		return C.ulong(x)
	default:
		return 65536
	}

}

//export SendInterrupt
func SendInterrupt(i C.ulong) {

	if py32.Subscribed(constants.Interrupt(i)) {
		py32SubbedInterruptChannel <- constants.Interrupt(i)
	}

}

//export GetRegister
func GetRegister(r C.ulong) C.ulong {

	return C.ulong(py32.Registers[r])

}

//export GetBuffer
func GetBuffer(b C.ulong) []C.char {

	//Convert to C.char array

	char_array := make([]C.char, 128)

	for index := range char_array {

		if constants.Register(b) == constants.RVideoText {
			char_array[index] = C.char(py32.TextBuffer[index])
		} else {
			char_array[index] = C.char(py32.DataBuffer[index])
		}

	}

	return char_array

}

//export SetRegister
func SetRegister(r C.ulong, v C.ulong) {

	py32.RegisterSync.Lock()
	py32.Registers[r] = uint32(v)
	py32.RegisterSync.Unlock()

}
