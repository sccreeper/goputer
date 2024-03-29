package compiler

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sccreeper/goputer/pkg/util"
	"time"
)

type AssembledProgram struct {
	ProgramBytes     []byte
	ProgramStructure ProgramStructure
	ProgramJson      string
}

// Compiler method
//
// Mainly used by the command line.
// Designed to work on desktop systems only.
func Compile(root_path string, get_file func(path string) []byte, config CompilerConfig, error_handler func(error_type ErrorType, error_text string)) (AssembledProgram, error) {

	start_time := time.Now().UnixMicro()

	prev_dir, err := os.Getwd()
	util.CheckError(err)

	//Change working directory so file imports are relative
	os.Chdir(filepath.Dir(config.FilePath))

	//Parse

	if config.Verbose {
		log.Println("Parsing...")
	}

	p := Parser{
		CodeString:   string(get_file(root_path)),
		FileName:     config.FilePath,
		Verbose:      false,
		Imported:     false,
		ErrorHandler: error_handler,
		FileReader:   get_file,
	}

	program_data, err := p.Parse()

	util.CheckError(err)

	// -----------------
	// Generate JSON
	// ----------------

	json_bytes, err := json.MarshalIndent(program_data, "", "\t")

	util.CheckError(err)

	// Begin bytecode generation

	if config.Verbose {
		log.Println("Bytecode generation...")
	}

	program_bytes := GenerateBytecode(program_data)

	//Output start indexes

	// log.Printf("Data start index: %d", data_start_index)
	// log.Printf("Jump start index: %d", jmp_block_start_index)
	// log.Printf("Interrupt table start index: %d", interrupt_table_start_index)
	// log.Printf("Program start index: %d", instruction_start_index)
	fmt.Printf("Final executable size: %d byte(s)\n", len(program_bytes))

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

	fmt.Printf("Compiled in %f %ss\n", elapsed_time, time_unit)

	os.Chdir(prev_dir)

	return AssembledProgram{

		ProgramBytes:     program_bytes,
		ProgramStructure: program_data,
		ProgramJson:      string(json_bytes),
	}, nil

}
