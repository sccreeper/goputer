package main

import (
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

		return ""

	})

	return init_func
}

func Run() js.Func {
	run_func := js.FuncOf(func(this js.Value, args []js.Value) any {

		go js32.Run()

		return ""

	})

	return run_func
}

// Methods for interacting with the VM

func SetRegister() js.Func {
	var set_register js.Func
	return set_register
}

func GetRegister() js.Func {
	var get_register js.Func
	return get_register
}

func GetBuffer() js.Func {
	var get_buffer js.Func
	return get_buffer
}

func SendInterrupt() js.Func {
	var send_interrupt js.Func
	return send_interrupt
}

func GetInterrupt() js.Func {
	var get_interrupt js.Func
	return get_interrupt
}

func IsSubscribed() js.Func {
	var is_subscribed js.Func
	return is_subscribed
}

func IsFinished() js.Func {
	var is_finished js.Func
	return is_finished
}

func main() {
	fmt.Println("GO WASM")

	js.Global().Set("compileCode", Compile())
	js.Global().Set("initVM", Init())
	js.Global().Set("runVM", Run())

	<-make(chan bool)

}
