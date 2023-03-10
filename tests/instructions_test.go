// Contains generalized test for all instructions that don't require special conditions.
package tests

import (
	"embed"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"
	"sccreeper/goputer/pkg/vm"
	"testing"

	"github.com/BurntSushi/toml"
)

//go:embed test_files
var test_files embed.FS

// Generalized instruction tests.

type TestArray struct {
	Tests []TestDetails `toml:"tests"`
}

type TestDetails struct {
	Name          string `toml:"name"`
	IsFile        bool   `toml:"is_file"`
	CheckValue    uint32 `toml:"check_value"`
	CheckRegister string `toml:"check_reg"`
	CodeText      string `toml:"text"`
}

func compile(text string) []byte {

	p := compiler.Parser{
		CodeString:   text,
		FileName:     "main.gpasm",
		Verbose:      false,
		Imported:     false,
		ErrorHandler: func(error_type compiler.ErrorType, error_text string) { panic(error_text) },
		FileReader:   func(path string) []byte { return []byte(text) },
	}

	program_structure, err := p.Parse()
	util.CheckError(err)

	program_bytes := compiler.GenerateBytecode(program_structure)

	return program_bytes

}

func TestInstructions(t *testing.T) {

	var test_details TestArray

	toml_file, err := test_files.ReadFile("test_files/instruction_tests.toml")
	if err != nil {
		panic(err)
	}

	toml.Unmarshal(toml_file, &test_details)

	// Start VM runtime

	for _, v := range test_details.Tests {

		// Compile example code

		program_bytes := compile(v.CodeText)

		// Create VM instance
		// TODO: make this more time and memory efficient.

		var test_32 vm.VM
		var test_32_interrupt_channel chan constants.Interrupt
		var test_32_subbed_interrupt_channel chan constants.Interrupt

		test_32_interrupt_channel = make(chan constants.Interrupt)
		test_32_subbed_interrupt_channel = make(chan constants.Interrupt)

		vm.InitVM(&test_32, program_bytes, test_32_interrupt_channel, test_32_subbed_interrupt_channel, false, false)

		test_32.Run()

		for {
			if !test_32.Finished {
				continue
			} else {
				break
			}
		}

		if test_32.Registers[constants.RegisterInts[v.CheckRegister]] != uint32(v.CheckValue) {
			t.Errorf("Failed instruction test %s. Value should be %d but got %d", v.Name, v.CheckValue, test_32.Registers[constants.RegisterInts[v.CheckRegister]])
		}

	}

}

// Basic instructions

func TestJump(t *testing.T) {

	var test_32 vm.VM
	var test_32_interrupt_channel chan constants.Interrupt
	var test_32_subbed_interrupt_channel chan constants.Interrupt

	test_32_interrupt_channel = make(chan constants.Interrupt)
	test_32_subbed_interrupt_channel = make(chan constants.Interrupt)

	program_text, err := test_files.ReadFile("test_files/test_jump.gpasm")
	if err != nil {
		panic(err)
	}

	vm.InitVM(&test_32, compile(string(program_text[:])), test_32_interrupt_channel, test_32_subbed_interrupt_channel, true, false)

	var in_jump bool = false
	var jump_addr uint32

	for {

		test_32.Step()

		if test_32.Opcode == constants.IJump && !in_jump {
			in_jump = true
			jump_addr = test_32.ArgLarge
		} else if test_32.Opcode != constants.IJump && in_jump {

			if test_32.Registers[constants.RProgramCounter]-5 != jump_addr {
				t.Fatalf("Program counter should be %d is %d instead\n", jump_addr, test_32.Registers[constants.RProgramCounter])
			} else {
				break
			}

		}

	}

}

func TestCall(t *testing.T) {
	var test_32 vm.VM
	var test_32_interrupt_channel chan constants.Interrupt
	var test_32_subbed_interrupt_channel chan constants.Interrupt

	test_32_interrupt_channel = make(chan constants.Interrupt)
	test_32_subbed_interrupt_channel = make(chan constants.Interrupt)

	program_text, err := test_files.ReadFile("test_files/test_jump.gpasm")
	if err != nil {
		panic(err)
	}

	vm.InitVM(&test_32, compile(string(program_text[:])), test_32_interrupt_channel, test_32_subbed_interrupt_channel, true, false)

	var in_call bool = false
	var call_addr uint32

	for {

		test_32.Step()

		if test_32.Opcode == constants.IJump && !in_call {
			in_call = true
			call_addr = test_32.ArgLarge
		} else if test_32.Opcode != constants.IJump && in_call {

			if test_32.Registers[constants.RProgramCounter]-5 != call_addr {
				t.Fatalf("Program counter should be %d is %d instead\n", call_addr, test_32.Registers[constants.RProgramCounter])
			} else {
				break
			}

		}

	}

}

// // Logical instructions

// Should jump, if it doesn't test fails

func TestConditionalJump(t *testing.T) {
	var test_32 vm.VM
	var test_32_interrupt_channel chan constants.Interrupt
	var test_32_subbed_interrupt_channel chan constants.Interrupt

	test_32_interrupt_channel = make(chan constants.Interrupt)
	test_32_subbed_interrupt_channel = make(chan constants.Interrupt)

	program_text, err := test_files.ReadFile("test_files/test_cndjump.gpasm")
	if err != nil {
		panic(err)
	}

	vm.InitVM(&test_32, compile(string(program_text[:])), test_32_interrupt_channel, test_32_subbed_interrupt_channel, true, false)

	var in_jump bool = false
	var jump_addr uint32

	for {

		test_32.Step()

		if test_32.Opcode == constants.IConditionalJump && !in_jump {
			in_jump = true
			jump_addr = test_32.ArgLarge
		} else if test_32.Opcode != constants.IConditionalJump && in_jump {

			if test_32.Registers[constants.RProgramCounter]-5 != jump_addr {
				t.Fatalf("Program counter should be %d is %d instead\n", jump_addr, test_32.Registers[constants.RProgramCounter])
			} else {
				break
			}

		}

	}

}

func TestConditionalCall(t *testing.T) {
	var test_32 vm.VM
	var test_32_interrupt_channel chan constants.Interrupt
	var test_32_subbed_interrupt_channel chan constants.Interrupt

	test_32_interrupt_channel = make(chan constants.Interrupt)
	test_32_subbed_interrupt_channel = make(chan constants.Interrupt)

	program_text, err := test_files.ReadFile("test_files/test_cndcall.gpasm")
	if err != nil {
		panic(err)
	}

	vm.InitVM(&test_32, compile(string(program_text[:])), test_32_interrupt_channel, test_32_subbed_interrupt_channel, true, false)

	var in_jump bool = false
	var jump_addr uint32

	for {

		test_32.Step()

		if test_32.Opcode == constants.IConditionalCall && !in_jump {
			in_jump = true
			jump_addr = test_32.ArgLarge
		} else if test_32.Opcode != constants.IConditionalCall && in_jump {

			if test_32.Registers[constants.RProgramCounter]-5 != jump_addr {
				t.Fatalf("Program counter should be %d is %d instead\n", jump_addr, test_32.Registers[constants.RProgramCounter])
			} else {
				break
			}

		}

	}
}

// func TestShiftLeft() {

// }

// func TestShiftRight() {

// }

// // Misc

// func TestHalt() {

// }

// func TestInterrupt() {

// }
