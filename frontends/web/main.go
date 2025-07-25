//go:build js && wasm

// WASM "proxy" layer between JS and goputer
package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/gpimg"
	"sccreeper/goputer/pkg/util"
	"sccreeper/goputer/pkg/vm"
	"syscall/js"

	_ "image/jpeg"
	_ "image/png"
)

var js32 vm.VM

var programBytes []byte

var itnMap map[uint32]string
var registerMap map[uint32]string
var interruptMap map[constants.Interrupt]string

type FileType string

const (
	textFile  FileType = "text"
	imageFile FileType = "image"
	binFile   FileType = "bin"
)

type File struct {
	Type    FileType
	Data    []byte
	Encoded bool
}

var fileMap map[string]File

//Custom compile and run methods.

func Compile() js.Func {

	compileFunc := js.FuncOf(func(this js.Value, args []js.Value) any {

		p := compiler.Parser{
			CodeString:   string(fileMap["main.gpasm"].Data),
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
		return []byte(val.Data), nil
	} else {
		return nil, compiler.ErrFile
	}
}

func UpdateFile(this js.Value, args []js.Value) any {

	fileData := make([]byte, args[2].Int())
	js.CopyBytesToGo(fileData, args[1])

	// Wether to encode/re-encode
	if args[4].Bool() && FileType(args[3].String()) == imageFile {

		bytesSource := bytes.NewReader(fileData)
		bytesDest := util.NewMemWriteSeeker()

		err := gpimg.Encode(
			bytesSource,
			bytesDest,
			gpimg.FlagRLECompression,
		)
		util.CheckError(err)

		fileData = bytesDest.Bytes()

	}

	fileMap[args[0].String()] = File{
		Data: fileData,
		Type: FileType(args[3].String()),
	}

	return js.ValueOf(nil)
}

func GetFile(this js.Value, args []js.Value) any {

	js.CopyBytesToJS(args[1], fileMap[args[0].String()].Data)

	return js.ValueOf(nil)

}

func GetFileSize(this js.Value, args []js.Value) any {
	return js.ValueOf(len(fileMap[args[0].String()].Data))
}

func DoesFileExist(this js.Value, args []js.Value) any {

	_, exists := fileMap[args[0].String()]

	return js.ValueOf(exists)

}

func GetFileType(this js.Value, args []js.Value) any {
	return js.ValueOf(string(fileMap[args[0].String()].Type))
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

func NumFiles(this js.Value, args []js.Value) any {
	return js.ValueOf(len(fileMap))
}

func HandleError(errorType compiler.ErrorMessage, errorText string) {

	js.Global().Call("showError", 1, js.ValueOf(errorText))

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

func GetRegisterBytes(this js.Value, args []js.Value) any {

	src := make([]byte, 4)
	binary.LittleEndian.PutUint32(src, js32.Registers[constants.Register(args[0].Int())])

	js.CopyBytesToJS(args[1], src)

	return js.ValueOf(nil)

}

func GetBuffer(this js.Value, args []js.Value) any {

	if args[0].String() == "text" {

		js.CopyBytesToJS(args[1], js32.TextBuffer[:])

		return js.ValueOf(nil)

	} else if args[0].String() == "data" {

		js.CopyBytesToJS(args[1], js32.DataBuffer[:])

		return js.ValueOf(nil)

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

func GetCurrentInstruction(this js.Value, args []js.Value) any {

	itn, err := compiler.DecodeInstructionString(js32.CurrentInstruction)

	if err != nil {
		return js.ValueOf(err.Error())
	} else {
		return js.ValueOf(itn)
	}

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

	return js.ValueOf(programBytes)

}

func SetProgramBytes(this js.Value, args []js.Value) any {

	programBytes = make([]byte, args[1].Int())

	js.CopyBytesToGo(programBytes, args[0])

	return js.ValueOf(nil)

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
	js.Global().Set("getRegisterBytes", js.FuncOf(GetRegisterBytes))
	js.Global().Set("getBuffer", js.FuncOf(GetBuffer))

	js.Global().Set("getInterrupt", js.FuncOf(GetInterrupt))
	js.Global().Set("sendInterrupt", js.FuncOf(SendInterrupt))
	js.Global().Set("isSubscribed", js.FuncOf(IsSubscribed))

	js.Global().Set("currentItn", js.FuncOf(GetCurrentInstruction))

	js.Global().Set("cycleVM", js.FuncOf(Cycle))

	js.Global().Set("isFinished", js.FuncOf(IsFinished))

	js.Global().Set("updateFile", js.FuncOf(UpdateFile))
	js.Global().Set("removeFile", js.FuncOf(RemoveFile))
	js.Global().Set("getFile", js.FuncOf(GetFile))
	js.Global().Set("getFiles", js.FuncOf(GetFiles))
	js.Global().Set("numFiles", js.FuncOf(NumFiles))
	js.Global().Set("getFileSize", js.FuncOf(GetFileSize))
	js.Global().Set("doesFileExist", js.FuncOf(DoesFileExist))
	js.Global().Set("getFileType", js.FuncOf(GetFileType))

	js.Global().Set("getProgramBytes", js.FuncOf(GetProgramBytes))
	js.Global().Set("setProgramBytes", js.FuncOf(SetProgramBytes))

	js.Global().Set("disassembleCode", js.FuncOf(Disassemble))

	js.Global().Set("updateFramebuffer", js.FuncOf(UpdateFrameBuffer))

	fileMap = make(map[string]File)
	fileMap["main.gpasm"] = File{
		Type: textFile,
		Data: []byte(""),
	}

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
	js.Global().Set("memSize", js.ValueOf(int(vm.MemSize-vm.VideoBufferSize))) // usable memory size

	js.Global().Set("convertColour", js.FuncOf(ConvertColour))
	js.Global().Set("convertHex", js.FuncOf(ConvertHex))

	<-make(chan bool)

}
