// WASM "proxy" layer between JS and goputer
package main

import (
	"errors"
	"fmt"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"
	"sccreeper/goputer/pkg/vm"
	"syscall/js"
)

var js32 vm.VM

var js32InterruptChannel chan constants.Interrupt = make(chan constants.Interrupt)
var js32SubbedInterruptChannel chan constants.Interrupt = make(chan constants.Interrupt)
var js32StepChannel chan bool = make(chan bool)

var program_bytes []byte

//Custom compile and run methods.

func Compile() js.Func {

	compile_func := js.FuncOf(func(this js.Value, args []js.Value) any {

		p := compiler.Parser{
			CodeString: args[0].String(),
			FileName:   "",
			Verbose:    false,
			Imported:   false,
		}

		program_structure, err := p.Parse()
		util.CheckError(err)

		program_bytes = compiler.GenerateBytecode(program_structure)

		fmt.Println(len(program_bytes))

		return ""

	})

	return compile_func

}

func Init() js.Func {

	init_func := js.FuncOf(func(this js.Value, args []js.Value) any {

		js32 = vm.VM{}

		vm.InitVM(&js32, program_bytes, js32.InterruptChannel, js32SubbedInterruptChannel, true, js32StepChannel)

		return js.Null

	})

	return init_func
}

func Run() js.Func {
	run_func := js.FuncOf(func(this js.Value, args []js.Value) any {

		go js32.Run()

		return js.Null

	})

	return run_func
}

// Methods for interacting with the VM

func SetRegister(this js.Value, args []js.Value) any {

	js32.RegisterSync.Lock()
	js32.Registers[constants.RegisterInts[args[0].String()]] = uint32(args[1].Int())
	js32.RegisterSync.Unlock()

	return js.Null

}

func GetRegister(this js.Value, args []js.Value) any {

	return js.ValueOf(js32.Registers[constants.RegisterInts[args[0].String()]])

}

func GetBuffer(this js.Value, args []js.Value) any {

	if args[0].String() == "text" {

		return js.ValueOf(js32.TextBuffer)

	} else if args[0].String() == "video" {

		return js.ValueOf(js32.TextBuffer)

	} else {

		panic(errors.New(fmt.Sprintf("'%s' is not a buffer.", args[0].String())))

	}

}

func SendInterrupt(this js.Value, args []js.Value) any {

	if js32.Subscribed(constants.Interrupt(constants.InterruptInts[args[0].String()])) {

		js32SubbedInterruptChannel <- constants.Interrupt(constants.InterruptInts[args[0].String()])

		return js.Null

	}

	return js.Null

}

func GetInterrupt(this js.Value, args []js.Value) any {

	select {
	case x := <-js32.InterruptChannel:
		return js.ValueOf(x)
	default:
		return js.Null

	}

}

func IsSubscribed(this js.Value, args []js.Value) any {

	return js.ValueOf(
		js32.Subscribed(
			constants.Interrupt(constants.InterruptInts[args[0].String()]),
		))

}

func IsFinished(this js.Value, args []js.Value) any {

	return js.ValueOf(js32.Finished)

}

func Step(this js.Value, args []js.Value) any {

	js32StepChannel <- true

	return js.Null

}

func main() {
	fmt.Println("JS32 init...")

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

	js.Global().Set("isFinished", js.FuncOf(IsFinished))

	<-make(chan bool)

}
