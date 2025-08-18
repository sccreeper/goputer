package main

/*
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"log"
	"math"
	"runtime"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/expansions"
	"sccreeper/goputer/pkg/vm"
	"unsafe"
)

var py32 *vm.VM
var videoBuffer *C.uint8_t

func main() {}

//export Init
func Init(programBytes *C.char, codeLength C.int) {

	py32, _ = vm.NewVM(
		C.GoBytes(unsafe.Pointer(programBytes), codeLength),
		expansions.ModuleExists,
		expansions.Interaction,
	)

	expansions.LoadExpansions(py32)

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

		py32.SubscribedInterruptQueue = append(py32.SubscribedInterruptQueue, constants.Interrupt(i))

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

//export GetInstructionString
func GetInstructionString() *C.char {

	itn, _ := compiler.DecodeInstructionString(py32.CurrentInstruction)
	return C.CString(itn)

}

//export GetVideoBufferPtr
func GetVideoBufferPtr() *C.uint8_t {

	videoBuffer = (*C.uint8_t)(C.malloc(C.size_t(320 * 240 * 3)))

	var pinner runtime.Pinner
	pinner.Pin(videoBuffer)

	C.memset(unsafe.Pointer(videoBuffer), 0, C.size_t(320*240*3))

	return videoBuffer

}

//export CopyFrameBuffer
func CopyFrameBuffer() {

	C.memcpy(unsafe.Pointer(videoBuffer), unsafe.Pointer(&py32.MemArray[0]), C.size_t(320*240*3))

}

//export SetExpansionModuleAttribute
func SetExpansionModuleAttribute(id *C.char, attrib *C.char, val *C.uint8_t, bytesLength C.int) {

	idStr := C.GoString(id)
	nameStr := C.GoString(attrib)
	valBytes := C.GoBytes(unsafe.Pointer(val), bytesLength)

	expansions.SetAttribute(idStr, nameStr, valBytes)
}

//export GetExpansionModuleAttribute
func GetExpansionModuleAttribute(id *C.char, attrib *C.char) *C.char {

	idStr := C.GoString(id)
	defer C.free(unsafe.Pointer(&idStr))

	attribString := C.GoString(id)
	defer C.free(unsafe.Pointer(&attribString))

	return C.CString(string(expansions.GetAttribute(idStr, attribString)))

}

// For Python on Windows

//export CFree
func CFree(ptr *C.uint8_t) {
	C.free(unsafe.Pointer(ptr))
}
