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
var testFiles embed.FS

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

	programStructure, err := p.Parse()
	util.CheckError(err)

	programBytes := compiler.GenerateBytecode(programStructure)

	return programBytes

}

func TestInstructions(t *testing.T) {

	var testDetails TestArray

	tomlFile, err := testFiles.ReadFile("test_files/instruction_tests.toml")
	if err != nil {
		panic(err)
	}

	toml.Unmarshal(tomlFile, &testDetails)

	// Start VM runtime

	for _, v := range testDetails.Tests {

		// Compile example code

		programBytes := compile(v.CodeText)

		// Create VM instance
		// TODO: make this more time and memory efficient.

		var test32 vm.VM

		vm.InitVM(&test32, programBytes, false)

		for {
			test32.Cycle()

			if !test32.Finished {
				continue
			} else {
				break
			}
		}

		if test32.Registers[constants.RegisterInts[v.CheckRegister]] != uint32(v.CheckValue) {
			t.Errorf("Failed instruction test %s. Value should be %d but got %d", v.Name, v.CheckValue, test32.Registers[constants.RegisterInts[v.CheckRegister]])
		}

	}

}

// Basic instructions

func TestJump(t *testing.T) {

	var test32 vm.VM

	programText, err := testFiles.ReadFile("test_files/test_jump.gpasm")
	if err != nil {
		panic(err)
	}

	vm.InitVM(&test32, compile(string(programText[:])), false)

	var inJump bool = false
	var jumpAddr uint32

	for {

		test32.Cycle()

		if test32.Opcode == constants.IJump && !inJump {
			inJump = true
			jumpAddr = test32.ArgLarge
		} else if test32.Opcode != constants.IJump && inJump {

			if test32.Registers[constants.RProgramCounter]-5 != jumpAddr {
				t.Fatalf("Program counter should be %d is %d instead\n", jumpAddr, test32.Registers[constants.RProgramCounter])
			} else {
				break
			}

		}

	}

}

func TestCall(t *testing.T) {

	var test32 vm.VM

	programText, err := testFiles.ReadFile("test_files/test_jump.gpasm")
	if err != nil {
		panic(err)
	}

	vm.InitVM(&test32, compile(string(programText[:])), false)

	var inCall bool = false
	var callAddr uint32

	for {

		test32.Cycle()

		if test32.Opcode == constants.IJump && !inCall {
			inCall = true
			callAddr = test32.ArgLarge
		} else if test32.Opcode != constants.IJump && inCall {

			if test32.Registers[constants.RProgramCounter]-5 != callAddr {
				t.Fatalf("Program counter should be %d is %d instead\n", callAddr, test32.Registers[constants.RProgramCounter])
			} else {
				break
			}

		}

	}

}

// // Logical instructions

// Should jump, if it doesn't test fails

func TestConditionalJump(t *testing.T) {

	var test32 vm.VM

	programText, err := testFiles.ReadFile("test_files/test_cndjump.gpasm")
	if err != nil {
		panic(err)
	}

	vm.InitVM(&test32, compile(string(programText[:])), false)

	var inJump bool = false
	var jumpAddr uint32

	for {

		test32.Cycle()

		if test32.Opcode == constants.IConditionalJump && !inJump {
			inJump = true
			jumpAddr = test32.ArgLarge
		} else if test32.Opcode != constants.IConditionalJump && inJump {

			if test32.Registers[constants.RProgramCounter]-5 != jumpAddr {
				t.Fatalf("Program counter should be %d is %d instead\n", jumpAddr, test32.Registers[constants.RProgramCounter])
			} else {
				break
			}

		}

	}

}

func TestConditionalCall(t *testing.T) {
	var test32 vm.VM

	programText, err := testFiles.ReadFile("test_files/test_cndcall.gpasm")
	if err != nil {
		panic(err)
	}

	vm.InitVM(&test32, compile(string(programText[:])), false)

	var inJump bool = false
	var jumpAddr uint32

	for {

		test32.Cycle()

		if test32.Opcode == constants.IConditionalCall && !inJump {
			inJump = true
			jumpAddr = test32.ArgLarge
		} else if test32.Opcode != constants.IConditionalCall && inJump {

			if test32.Registers[constants.RProgramCounter]-5 != jumpAddr {
				t.Fatalf("Program counter should be %d is %d instead\n", jumpAddr, test32.Registers[constants.RProgramCounter])
			} else {
				break
			}

		}

	}
}
