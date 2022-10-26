package compiler

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"sccreeper/govm/pkg/util"
	"time"
)

// Assembler method
func Assemble(code_string string, config CompilerConfig) error {

	start_time := time.Now().UnixMicro()

	//Parse

	program_data, err := parse(code_string)

	util.CheckError(err)

	// Begin bytecode generation

	log.Println("Starting bytecode generation")

	final_byte_array := generate_bytecode(program_data)

	//Write to file

	os.WriteFile(config.OutputPath, final_byte_array, 0666)
	//Output start indexes

	// log.Printf("Data start index: %d", data_start_index)
	// log.Printf("Jump start index: %d", jmp_block_start_index)
	// log.Printf("Interrupt table start index: %d", interrupt_table_start_index)
	// log.Printf("Program start index: %d", instruction_start_index)
	log.Printf("Final executable size: %d byte(s)", len(final_byte_array))

	// -----------------
	// Output JSON
	// ----------------

	if config.OutputJSON {

		json_bytes, err := json.MarshalIndent(program_data, "", "\t")

		util.CheckError(err)

		err = os.WriteFile(config.JSONPath, json_bytes, 0666)

		util.CheckError(err)

		log.Printf("Outputted JSON structure to '%s'", config.JSONPath)

	}

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

	return nil

}
