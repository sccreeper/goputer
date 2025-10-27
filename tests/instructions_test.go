// Contains generalized test for all instructions that don't require special conditions.
package tests

import (
	"embed"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/expansions"
	"sccreeper/goputer/pkg/util"
	"sccreeper/goputer/pkg/vm"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
)

//go:embed test_files/instructions
var instructionTestFiles embed.FS

// Generalized instruction tests.

type TestArray struct {
	Tests []TestDetails `toml:"tests"`
}

type TestDetails struct {
	Name          string   `toml:"name"`
	IsFile        bool     `toml:"is_file"`
	CheckValue    []uint32 `toml:"check_value"`
	CheckRegister string   `toml:"check_reg"`
	CodeText      string   `toml:"text"`
}

func compile(text string) []byte {

	p := compiler.Parser{
		CodeString:   text,
		FileName:     "main.gpasm",
		Verbose:      false,
		Imported:     false,
		ErrorHandler: func(error_type compiler.ErrorMessage, error_text string) { panic(error_text) },
		FileReader:   func(path string) ([]byte, error) { return []byte(text), nil },
	}

	programStructure, err := p.Parse()
	util.CheckError(err)

	programBytes := compiler.GenerateBytecode(programStructure, false)

	return programBytes

}

func TestInstructions(t *testing.T) {

	var testDetails TestArray

	tomlFile, err := instructionTestFiles.ReadFile("test_files/instructions/instruction_tests.toml")
	if err != nil {
		panic(err)
	}

	toml.Unmarshal(tomlFile, &testDetails)

	// Start VM runtime

	for _, v := range testDetails.Tests {

		t.Run(v.Name, func(t *testing.T) {
			// Compile example code

			programBytes := compile(v.CodeText)

			// Create VM instance
			// TODO: make this more time and memory efficient.

			test32, _ := vm.NewVM(programBytes, expansions.ModuleExists, expansions.Interaction)

			for {
				if test32.Finished {
					break
				}

				test32.Cycle()
			}

			if len(v.CheckValue) == 1 {
				if test32.Registers[constants.RegisterInts[v.CheckRegister]] != uint32(v.CheckValue[0]) {
					t.Errorf("Failed instruction test %s. Value should be %d but got %d\n", v.Name, v.CheckValue[0], test32.Registers[constants.RegisterInts[v.CheckRegister]])
				}
			} else if len(v.CheckValue) == 2 {
				if !(test32.Registers[constants.RegisterInts[v.CheckRegister]] <= v.CheckValue[1] && test32.Registers[constants.RegisterInts[v.CheckRegister]] >= v.CheckValue[0]) {
					t.Errorf("Failed instruction test %s. Value should be between %d and %d (inclusive) but got %d\n", v.Name, v.CheckValue[0], v.CheckValue[1], test32.Registers[constants.RegisterInts[v.CheckRegister]])
				}
			}
		})

	}

}

func BenchmarkInstructions(b *testing.B) {

	var benchmarkDetails TestArray

	tomlFile, err := instructionTestFiles.ReadFile("test_files/instructions/instruction_tests.toml")
	if err != nil {
		panic(err)
	}

	toml.Unmarshal(tomlFile, &benchmarkDetails)

	// Start VM runtime

	for _, v := range benchmarkDetails.Tests {

		programBytes := compile(v.CodeText)
		test32, _ := vm.NewVM(programBytes, expansions.ModuleExists, expansions.Interaction)

		initialProgramCounter := test32.Registers[constants.RProgramCounter]

		b.Run(v.Name, func(b *testing.B) {

			for b.Loop() {

				for {
					if test32.Finished {
						break
					}

					test32.Cycle()
				}

				test32.Finished = false
				test32.Registers[constants.RProgramCounter] = initialProgramCounter

			}

		})

	}

}

func BenchmarkTimeNow(b *testing.B) {

	for b.Loop() {

		time.Now()

	}

}

// func BenchmarkKeyChecking(b *testing.B) {

// 	b.Run("map_test", func(b *testing.B) {

// 		testMap := make(map[int]string, 0)

// 		for b.Loop() {

// 			if _, exists := testMap[b.N]; !exists {
// 				testMap[b.N] = "hello world"
// 			}

// 		}

// 	})

// }

// Basic instructions

func TestJump(t *testing.T) {

	programText, err := instructionTestFiles.ReadFile("test_files/instructions/test_jump.gpasm")
	if err != nil {
		panic(err)
	}

	test32, _ := vm.NewVM(compile(string(programText[:])), expansions.ModuleExists, expansions.Interaction)

	var inJump bool = false
	var jumpAddr uint32

	for {

		if test32.Finished {
			break
		}

		test32.Cycle()

		if test32.Opcode == constants.IJump && !inJump {
			inJump = true
			jumpAddr = test32.LongArg
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

	programText, err := instructionTestFiles.ReadFile("test_files/instructions/test_call.gpasm")
	if err != nil {
		panic(err)
	}

	test32, _ := vm.NewVM(compile(string(programText[:])), expansions.ModuleExists, expansions.Interaction)

	var inCall bool = false
	var callAddr uint32

	for {

		if test32.Finished {
			break
		}

		test32.Cycle()

		if test32.Opcode == constants.ICall && !inCall {
			inCall = true
			callAddr = test32.LongArg
		} else if test32.Opcode != constants.ICall && inCall {

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

	programText, err := instructionTestFiles.ReadFile("test_files/instructions/test_cndjump.gpasm")
	if err != nil {
		panic(err)
	}

	test32, _ := vm.NewVM(compile(string(programText[:])), expansions.ModuleExists, expansions.Interaction)

	var inJump bool = false
	var jumpAddr uint32

	for {

		if test32.Finished {
			break
		}

		test32.Cycle()

		if test32.Opcode == constants.IConditionalJump && !inJump {
			inJump = true
			jumpAddr = test32.LongArg
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

	programText, err := instructionTestFiles.ReadFile("test_files/instructions/test_cndcall.gpasm")
	if err != nil {
		panic(err)
	}

	test32, _ := vm.NewVM(compile(string(programText[:])), expansions.ModuleExists, expansions.Interaction)

	var inCall bool = false
	var jumpAddr uint32

	for {

		if test32.Finished {
			break
		}

		test32.Cycle()

		if test32.Opcode == constants.IConditionalCall && !inCall {
			inCall = true
			jumpAddr = test32.LongArg
		} else if test32.Opcode != constants.IConditionalCall && inCall {

			if test32.Registers[constants.RProgramCounter]-5 != jumpAddr {
				t.Fatalf("Program counter should be %d is %d instead\n", jumpAddr, test32.Registers[constants.RProgramCounter])
			} else {
				break
			}

		}

	}
}
