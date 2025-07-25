package rendering

import (
	"fmt"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"
	"sccreeper/goputer/pkg/vm"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var itnMap map[uint32]string
var registerMap map[uint32]string
var interuptMap map[constants.Interrupt]string

// Just reverses maps for now
func InitVMDebug() {

	itnMap = make(map[uint32]string)

	for k, v := range constants.InstructionInts {
		itnMap[v] = k
	}

	registerMap = make(map[uint32]string)

	for k, v := range constants.RegisterInts {
		registerMap[v] = k
	}

	interuptMap = make(map[constants.Interrupt]string)

	for k, v := range constants.InterruptInts {
		interuptMap[v] = k
	}

}

func RenderVMDebug(m *vm.VM) {

	rl.ClearBackground(rl.LightGray)

	rl.DrawText(fmt.Sprintf("Program counter: 0x%s", util.ConvertHex(m.Registers[constants.RProgramCounter])), 0, 0, 16, rl.Black)

	itn, err := compiler.DecodeInstructionString(m.CurrentInstruction)
	if err != nil {
		itn = err.Error()
	}

	rl.DrawText(fmt.Sprintf("Instruction: %s", itn), 0, 33, 16, rl.Black)

	// Buttons

}
