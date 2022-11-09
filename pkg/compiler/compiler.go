package compiler

import (
	"encoding/json"
	"log"
	"math"
	"sccreeper/govm/pkg/util"
	"time"
)

type AssembledProgram struct {
	ProgramBytes     []byte
	ProgramStructure ProgramStructure
	ProgramJson      string
}

// Assembler method
func Assemble(code_string string, config CompilerConfig) (AssembledProgram, error) {

	start_time := time.Now().UnixMicro()

	//Parse

	if config.Verbose {
		log.Println("Parsing...")
	}

	program_data, err := parse(code_string, config.Verbose)

	util.CheckError(err)

	// Begin bytecode generation

	if config.Verbose {
		log.Println("Bytecode generation...")
	}

	program_bytes := generate_bytecode(program_data)

	//Output start indexes

	// log.Printf("Data start index: %d", data_start_index)
	// log.Printf("Jump start index: %d", jmp_block_start_index)
	// log.Printf("Interrupt table start index: %d", interrupt_table_start_index)
	// log.Printf("Program start index: %d", instruction_start_index)
	log.Printf("Final executable size: %d byte(s)", len(program_bytes))

	// -----------------
	// Generate JSON
	// ----------------

	json_bytes, err := json.MarshalIndent(program_data, "", "\t")

	util.CheckError(err)

	// -------------------
	// Output elapsed time
	// -------------------

	elapsed_time := float64(time.Now().UnixMicro() - start_time)
	time_unit := "µ"

	if elapsed_time > math.Pow10(6) {
		elapsed_time = elapsed_time / math.Pow10(6)
		time_unit = ""
	} else if elapsed_time > math.Pow10(3) {
		elapsed_time = elapsed_time / math.Pow10(3)
		time_unit = "m"

	}

	log.Printf("Compiled in %f %ss", elapsed_time, time_unit)

	return AssembledProgram{

		ProgramBytes:     program_bytes,
		ProgramStructure: program_data,
		ProgramJson:      string(json_bytes),
	}, nil

}
