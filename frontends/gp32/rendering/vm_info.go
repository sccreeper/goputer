package rendering

import (
	"fmt"
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

	//Program counter as Hex

	rl.DrawText(fmt.Sprintf("Program counter: 0x%s", util.ConvertHex(m.Registers[constants.RProgramCounter])), 0, 0, 16, rl.Black)

	//Generate current instruction string

	var argText string = ""

	switch m.Opcode {
	case constants.IJump, constants.ICall, constants.IConditionalJump, constants.IConditionalCall:
		argText = util.ConvertHex(m.ArgLarge)
	default:
		if constants.InstructionArgumentCounts[m.Opcode][0] == 1 && m.Opcode != constants.ICallInterrupt {
			argText = registerMap[m.ArgLarge]
		} else if constants.InstructionArgumentCounts[m.Opcode][0] == 0 {
			argText = ""
		} else if m.Opcode == constants.ICallInterrupt {
			argText = interuptMap[constants.Interrupt(m.ArgLarge)]
		} else {
			argText = fmt.Sprintf("%s %s", registerMap[uint32(m.ArgSmall0)], registerMap[uint32(m.ArgSmall1)])
		}
	}

	rl.DrawText(fmt.Sprintf("Instruction: %s %s", itnMap[uint32(m.Opcode)], argText), 0, 33, 16, rl.Black)

}
