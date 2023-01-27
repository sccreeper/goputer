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

		p := compiler.Parser{
			CodeString:   v.CodeText,
			FileName:     "main.gpasm",
			Verbose:      false,
			Imported:     false,
			ErrorHandler: func(error_type compiler.ErrorType, error_text string) { t.Fatalf(error_text) },
			FileReader:   func(path string) []byte { return []byte(v.CodeText) },
		}

		program_structure, err := p.Parse()
		util.CheckError(err)

		program_bytes := compiler.GenerateBytecode(program_structure)

		// Create VM instance
		// TODO: make this more time and memory efficient.

		var test_32 vm.VM
		var test_32_interrupt_channel chan constants.Interrupt
		var test_32_subbed_interrupt_channel chan constants.Interrupt

		test_32_interrupt_channel = make(chan constants.Interrupt)
		test_32_subbed_interrupt_channel = make(chan constants.Interrupt)

		vm.InitVM(&test_32, program_bytes, test_32_interrupt_channel, test_32_subbed_interrupt_channel, false)

		test_32.Run()

		for {
			if !test_32.Finished {
				continue
			} else {
				break
			}
		}

		if test_32.Registers[constants.RegisterInts[v.CheckRegister]] != uint32(v.CheckValue) {
			t.Errorf("Failed instruction test %s", v.Name)
		}

	}

}

// Basic instructions

// func TestJump() {

// }

// func TestCall() {

// }

// // Logical instructions

// func TestConditionalJump() {

// }

// func TestConditionalCall() {

// }

// func TestOr() {

// }

// func TestXor() {

// }

// func TestAnd() {

// }

// func TestInvert() {

// }

// func TestShiftLeft() {

// }

// func TestShiftRight() {

// }

// func TestShiftEqual() {

// }

// func TestGreaterThan() {

// }

// func TestLessThan() {

// }

// // Misc

// func TestHalt() {

// }

// func TestInterrupt() {

// }

// func TestClear() {

// }
