package rendering

import (
	"fmt"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"
	"sccreeper/goputer/pkg/vm"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/slices"
)

var itn_map map[uint32]string
var register_map map[uint32]string
var interrupt_map map[constants.Interrupt]string

// Just reverses maps for now
func InitVMDebug() {

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

}

func RenderVMDebug(m *vm.VM) {

	rl.ClearBackground(rl.LightGray)

	//Program counter as Hex

	rl.DrawText(fmt.Sprintf("Program counter: 0x%s", util.ConvertHex(m.Registers[constants.RProgramCounter])), 0, 0, 16, rl.Black)

	//Generate current instruction string

	var arg_text string = ""

	switch m.Opcode {
	case constants.IJump, constants.ICall, constants.IConditionalJump, constants.IConditionalCall:
		arg_text = util.ConvertHex(m.ArgLarge)
	default:
		if slices.Contains(constants.SingleArgInstructions, m.Opcode) && m.Opcode != constants.ICallInterrupt {
			arg_text = register_map[m.ArgLarge]
		} else if m.Opcode == constants.ICallInterrupt {
			arg_text = interrupt_map[constants.Interrupt(m.ArgLarge)]
		} else {
			arg_text = fmt.Sprintf("%s %s", register_map[uint32(m.ArgSmall0)], register_map[uint32(m.ArgSmall1)])
		}
	}

	rl.DrawText(fmt.Sprintf("Instruction: %s %s", itn_map[uint32(m.Opcode)], arg_text), 0, 33, 16, rl.Black)

}
