//go:build js && wasm

// WASM "proxy" layer between JS and goputer
package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"
	"sccreeper/goputer/pkg/vm"
	"syscall/js"
)

var js32 vm.VM

var programBytes []byte

var itnMap map[uint32]string
var registerMap map[uint32]string
var interruptMap map[constants.Interrupt]string

var fileMap map[string]string

//Custom compile and run methods.

func Compile() js.Func {

	compileFunc := js.FuncOf(func(this js.Value, args []js.Value) any {

		p := compiler.Parser{
			CodeString:   args[0].String(),
			FileName:     "main.gpasm",
			Verbose:      false,
			Imported:     false,
			ErrorHandler: HandleError,
			FileReader:   fileReader,
		}

		programStructure, err := p.Parse()
		util.CheckError(err)

		programBytes = compiler.GenerateBytecode(programStructure, false)

		fmt.Printf("Compiled program length: %d\n", len(programBytes))

		return js.ValueOf(nil)

	})

	return compileFunc

}

func fileReader(path string) ([]byte, error) {
	if val, exists := fileMap[path]; exists {
		return []byte(val), nil
	} else {
		return nil, compiler.ErrFile
	}
}

func UpdateFile(this js.Value, args []js.Value) any {
	fileMap[args[0].String()] = args[1].String()

	return js.ValueOf(nil)
}

func GetFile(this js.Value, args []js.Value) any {

	return js.ValueOf(fileMap[args[0].String()])

}

func RemoveFile(this js.Value, args []js.Value) any {

	delete(fileMap, args[0].String())

	return js.ValueOf(nil)

}

func GetFiles(this js.Value, args []js.Value) any {

	keys := make([]interface{}, 0)

	for k, _ := range fileMap {
		keys = append(keys, k)
	}

	return js.ValueOf(keys)

}

func HandleError(error_type compiler.ErrorType, error_text string) {

	js.Global().Call("showError", 1, js.ValueOf(error_text))

}

func Init() js.Func {

	init_func := js.FuncOf(func(this js.Value, args []js.Value) any {

		js32 = vm.VM{}

		vm.InitVM(&js32, programBytes, true)

		return js.ValueOf(nil)

	})

	return init_func
}

// Methods for interacting with the VM

func SetRegister(this js.Value, args []js.Value) any {

	js32.Registers[constants.Register(args[0].Int())] = uint32(args[1].Int())

	return js.ValueOf(nil)

}

func GetRegister(this js.Value, args []js.Value) any {

	return js.ValueOf(js32.Registers[constants.Register(args[0].Int())])

}

func GetBuffer(this js.Value, args []js.Value) any {

	if args[0].String() == "text" {

		//Convert buffer

		converted := make([]interface{}, 0)

		for _, v := range js32.TextBuffer {

			converted = append(converted, int(v))

		}

		return js.ValueOf(converted)

	} else if args[0].String() == "data" {

		//Convert buffer

		converted := make([]interface{}, 0)

		for _, v := range js32.DataBuffer {

			converted = append(converted, int(v))

		}

		return js.ValueOf(converted)

	} else {

		panic(errors.New(fmt.Sprintf("'%s' is not a buffer.", args[0].String())))

	}

}

func SendInterrupt(this js.Value, args []js.Value) any {

	if js32.Subscribed(constants.Interrupt(args[0].Int())) {

		js32.SubbedInterruptQueue = append(js32.SubbedInterruptQueue, constants.Interrupt(args[0].Int()))

		return js.ValueOf(nil)

	}

	return js.ValueOf(nil)

}

func GetInterrupt(this js.Value, args []js.Value) any {

	if len(js32.InterruptQueue) > 0 {

		var x constants.Interrupt
		x, js32.InterruptQueue = js32.InterruptQueue[0], js32.InterruptQueue[1:]

		return js.ValueOf(int(x))

	} else {
		return js.ValueOf(nil)
	}

}

func IsSubscribed(this js.Value, args []js.Value) any {

	return js.ValueOf(
		js32.Subscribed(
			constants.Interrupt(args[0].Int()),
		))

}

func IsFinished(this js.Value, args []js.Value) any {

	return js.ValueOf(js32.Finished)

}

func UpdateFrameBuffer(this js.Value, args []js.Value) any {

	js.CopyBytesToJS(js.Global().Get("textureData"), js32.MemArray[:vm.VideoBufferSize])

	return js.ValueOf(nil)

}

func Cycle(this js.Value, args []js.Value) any {

	js32.Cycle()

	return js.ValueOf(nil)

}

func ParserItnStr(this js.Value, args []js.Value) any {

	//Generate current instruction string

	var arg_text string = ""

	switch js32.Opcode {
	case constants.IJump, constants.ICall, constants.IConditionalJump, constants.IConditionalCall:
		arg_text = util.ConvertHex(js32.ArgLarge)
	default:
		if constants.InstructionArgumentCounts[js32.Opcode][0] == 1 && js32.Opcode != constants.ICallInterrupt {
			arg_text = registerMap[js32.ArgLarge]
		} else if js32.Opcode == constants.ICallInterrupt {
			arg_text = interruptMap[constants.Interrupt(js32.ArgLarge)]
		} else {
			arg_text = fmt.Sprintf("%s %s", registerMap[uint32(js32.ArgSmall0)], registerMap[uint32(js32.ArgSmall1)])
		}
	}

	return js.ValueOf(fmt.Sprintf("%s %s", itnMap[uint32(js32.Opcode)], arg_text))

}

func Disassemble(this js.Value, args []js.Value) any {

	programIntArray := make([]byte, 0)

	for i := 0; i < args[0].Length(); i++ {

		programIntArray = append(programIntArray, byte(args[0].Index(i).Int()))

	}

	disassembledProgram, err := compiler.Disassemble(programIntArray, false)
	util.CheckError(err)

	jsonString, err := json.Marshal(disassembledProgram)
	util.CheckError(err)

	return js.ValueOf(string(jsonString[:]))

}

func GetProgramBytes(this js.Value, args []js.Value) any {

	// Convert to []interface{}

	interface_program_bytes := make([]interface{}, 0)

	for _, v := range programBytes {
		interface_program_bytes = append(interface_program_bytes, v)
	}

	return js.ValueOf(interface_program_bytes)

}

//Other

func ConvertColour(this js.Value, args []js.Value) any {

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b[:], uint32(args[0].Int()))

	return js.ValueOf(fmt.Sprintf("rgba(%d, %d, %d, %f)", b[0], b[1], b[2], math.Round(float64(b[3])/255)))

}

func ConvertHex(this js.Value, args []js.Value) any {

	if args[1].Bool() {
		return js.ValueOf(util.ConvertHex[int](args[0].Int() + int(compiler.StackSize)))
	} else {
		return js.ValueOf(util.ConvertHex[int](args[0].Int()))
	}

}

func main() {
	fmt.Println("JS32 init...")

	// Reversed maps

	itnMap = make(map[uint32]string)

	for k, v := range constants.InstructionInts {
		itnMap[v] = k
	}

	registerMap = make(map[uint32]string)

	for k, v := range constants.RegisterInts {
		registerMap[v] = k
	}

	interruptMap = make(map[constants.Interrupt]string)

	for k, v := range constants.InterruptInts {
		interruptMap[v] = k
	}

	// VM init methods

	js.Global().Set("compileCode", Compile())
	js.Global().Set("initVM", Init())

	//VM interaction methods

	js.Global().Set("setRegister", js.FuncOf(SetRegister))
	js.Global().Set("getRegister", js.FuncOf(GetRegister))
	js.Global().Set("getBuffer", js.FuncOf(GetBuffer))

	js.Global().Set("getInterrupt", js.FuncOf(GetInterrupt))
	js.Global().Set("sendInterrupt", js.FuncOf(SendInterrupt))
	js.Global().Set("isSubscribed", js.FuncOf(IsSubscribed))

	js.Global().Set("currentItn", js.FuncOf(ParserItnStr))

	js.Global().Set("cycleVM", js.FuncOf(Cycle))

	js.Global().Set("isFinished", js.FuncOf(IsFinished))

	js.Global().Set("updateFile", js.FuncOf(UpdateFile))
	js.Global().Set("removeFile", js.FuncOf(RemoveFile))
	js.Global().Set("getFile", js.FuncOf(GetFile))
	js.Global().Set("getFiles", js.FuncOf(GetFiles))

	js.Global().Set("getProgramBytes", js.FuncOf(GetProgramBytes))

	js.Global().Set("disassembleCode", js.FuncOf(Disassemble))

	js.Global().Set("updateFramebuffer", js.FuncOf(UpdateFrameBuffer))

	fileMap = make(map[string]string)
	fileMap["main.gpasm"] = ""

	//Convert constants maps into [string]interface maps

	interruptsConverted := make(map[string]interface{}, 0)

	for k, v := range constants.InterruptInts {
		interruptsConverted[k] = int(v)
	}

	instructionsConverted := make(map[string]interface{}, 0)

	for k, v := range constants.InstructionInts {
		instructionsConverted[k] = int(v)
	}

	registersConverted := make(map[string]interface{}, 0)

	for k, v := range constants.RegisterInts {
		registersConverted[k] = int(v)
	}

	// Make an instruction & interrupt array for disassembling

	instructionsArray := make([]interface{}, vm.InstructionCount)

	for k, v := range constants.InstructionInts {

		instructionsArray[v] = k

	}

	interruptArray := make([]interface{}, vm.InterruptCount)

	for k, v := range constants.InterruptInts {

		interruptArray[v] = k

	}

	js.Global().Set("interruptInts", js.ValueOf(interruptsConverted))
	js.Global().Set("instructionInts", js.ValueOf(instructionsConverted))
	js.Global().Set("registerInts", js.ValueOf(registersConverted))

	js.Global().Set("instructionArray", js.ValueOf(instructionsArray))
	js.Global().Set("interruptArray", js.ValueOf(interruptArray))

	js.Global().Set("memOffset", js.ValueOf(int(compiler.StackSize)))

	js.Global().Set("convertColour", js.FuncOf(ConvertColour))
	js.Global().Set("convertHex", js.FuncOf(ConvertHex))

	<-make(chan bool)

}
