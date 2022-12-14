// WASM "proxy" layer between JS and goputer
package main

import (
	"encoding/binary"
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

var js32InterruptChannel chan constants.Interrupt = make(chan constants.Interrupt)
var js32SubbedInterruptChannel chan constants.Interrupt = make(chan constants.Interrupt)

var program_bytes []byte

var itn_map map[uint32]string
var register_map map[uint32]string
var interrupt_map map[constants.Interrupt]string

//Custom compile and run methods.

func Compile() js.Func {

	compile_func := js.FuncOf(func(this js.Value, args []js.Value) any {

		p := compiler.Parser{
			CodeString:   args[0].String(),
			FileName:     "",
			Verbose:      false,
			Imported:     false,
			ErrorHandler: HandleError,
		}

		program_structure, err := p.Parse()
		util.CheckError(err)

		program_bytes = compiler.GenerateBytecode(program_structure)

		fmt.Printf("Compiled program length: %d\n", len(program_bytes))

		return js.ValueOf(nil)

	})

	return compile_func

}

func HandleError(error_type compiler.ErrorType, error_text string) {

	js.Global().Call("showError", 1, js.ValueOf(error_text))

}

func Init() js.Func {

	init_func := js.FuncOf(func(this js.Value, args []js.Value) any {

		js32 = vm.VM{}

		vm.InitVM(&js32, program_bytes, js32.InterruptChannel, js32SubbedInterruptChannel, true)

		return js.ValueOf(nil)

	})

	return init_func
}

func Run() js.Func {
	run_func := js.FuncOf(func(this js.Value, args []js.Value) any {

		go js32.Run()

		return js.ValueOf(nil)

	})

	return run_func
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

		js32.SubbedInterruptArray = append(js32.SubbedInterruptArray, constants.Interrupt(args[0].Int()))

		return js.ValueOf(nil)

	}

	return js.ValueOf(nil)

}

func GetInterrupt(this js.Value, args []js.Value) any {

	if len(js32.InterruptArray) > 0 {

		x := js32.InterruptArray[len(js32.InterruptArray)-1]
		js32.InterruptArray = js32.InterruptArray[:len(js32.InterruptArray)-1]

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

func Step(this js.Value, args []js.Value) any {

	js32.Step()

	return js.ValueOf(nil)

}

func ItnStr(this js.Value, args []js.Value) any {

	//Generate current instruction string

	var arg_text string = ""

	switch js32.Opcode {
	case constants.IJump, constants.ICall, constants.IConditionalJump, constants.IConditionalCall:
		arg_text = util.ConvertHex(js32.ArgLarge)
	default:
		if util.SliceContains(constants.SingleArgInstructions, js32.Opcode) && js32.Opcode != constants.ICallInterrupt {
			arg_text = register_map[js32.ArgLarge]
		} else if js32.Opcode == constants.ICallInterrupt {
			arg_text = interrupt_map[constants.Interrupt(js32.ArgLarge)]
		} else {
			arg_text = fmt.Sprintf("%s %s", register_map[uint32(js32.ArgSmall0)], register_map[uint32(js32.ArgSmall1)])
		}
	}

	return js.ValueOf(fmt.Sprintf("%s %s", itn_map[uint32(js32.Opcode)], arg_text))

}

//Other

func ConvertColour(this js.Value, args []js.Value) any {

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b[:], uint32(args[0].Int()))

	return js.ValueOf(fmt.Sprintf("rgba(%d, %d, %d, %f)", b[0], b[1], b[2], math.Round(float64(b[3])/255)))

}

func main() {
	fmt.Println("JS32 init...")

	// Reversed maps

	itn_map = make(map[uint32]string)

	for k, v := range constants.InstructionInts {
		itn_map[v] = k
	}

	register_map = make(map[uint32]string)

	for k, v := range constants.RegisterInts {
		register_map[v] = k
	}

	interrupt_map = make(map[constants.Interrupt]string)

	for k, v := range constants.InterruptInts {
		interrupt_map[v] = k
	}

	// VM init methods

	js.Global().Set("compileCode", Compile())
	js.Global().Set("initVM", Init())
	js.Global().Set("runVM", Run())

	//VM interaction methods

	js.Global().Set("setRegister", js.FuncOf(SetRegister))
	js.Global().Set("getRegister", js.FuncOf(GetRegister))
	js.Global().Set("getBuffer", js.FuncOf(GetBuffer))

	js.Global().Set("getInterrupt", js.FuncOf(GetInterrupt))
	js.Global().Set("sendInterrupt", js.FuncOf(SendInterrupt))
	js.Global().Set("isSubscribed", js.FuncOf(IsSubscribed))

	js.Global().Set("currentItn", js.FuncOf(ItnStr))

	js.Global().Set("stepVM", js.FuncOf(Step))

	js.Global().Set("isFinished", js.FuncOf(IsFinished))

	//Convert constants maps into [string]interface maps

	interrupts_converted := make(map[string]interface{}, 0)

	for k, v := range constants.InterruptInts {
		interrupts_converted[k] = int(v)
	}

	instructions_converted := make(map[string]interface{}, 0)

	for k, v := range constants.InstructionInts {
		instructions_converted[k] = int(v)
	}

	registers_converted := make(map[string]interface{}, 0)

	for k, v := range constants.RegisterInts {
		registers_converted[k] = int(v)
	}

	js.Global().Set("interruptInts", js.ValueOf(interrupts_converted))
	js.Global().Set("instructionInts", js.ValueOf(instructions_converted))
	js.Global().Set("registerInts", js.ValueOf(registers_converted))

	js.Global().Set("convertColour", js.FuncOf(ConvertColour))

	<-make(chan bool)

}
