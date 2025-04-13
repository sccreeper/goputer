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
func Compile(rootPath string, getFile func(path string) ([]byte, error), config CompilerConfig, error_handler func(errorType ErrorType, errorText string)) (AssembledProgram, error) {

	startTime := time.Now().UnixMicro()

	prevDir, err := os.Getwd()
	util.CheckError(err)

	//Change working directory so file imports are relative
	os.Chdir(filepath.Dir(config.FilePath))

	//Parse

	if config.Verbose {
		log.Println("Parsing...")
	}

	fileData, err := getFile(rootPath)

	if err != nil {
		return AssembledProgram{}, err
	}

	p := Parser{
		CodeString:   string(fileData),
		FileName:     config.FilePath,
		Verbose:      false,
		Imported:     false,
		ErrorHandler: error_handler,
		FileReader:   getFile,
	}

	programData, err := p.Parse()

	util.CheckError(err)

	// -----------------
	// Generate JSON
	// ----------------

	jsonBytes, err := json.MarshalIndent(programData, "", "\t")

	util.CheckError(err)

	// Begin bytecode generation

	if config.Verbose {
		log.Println("Bytecode generation...")
	}

	programBytes := GenerateBytecode(programData)

	//Output start indexes

	// log.Printf("Data start index: %d", data_start_index)
	// log.Printf("Jump start index: %d", jmp_block_start_index)
	// log.Printf("Interrupt table start index: %d", interrupt_table_start_index)
	// log.Printf("Program start index: %d", instruction_start_index)
	fmt.Printf("Final executable size: %d byte(s)\n", len(programBytes))

	// -------------------
	// Output elapsed time
	// -------------------

	elapsedTime := float64(time.Now().UnixMicro() - startTime)
	timeUnit := "µ"

	if elapsedTime > math.Pow10(6) {
		elapsedTime = elapsedTime / math.Pow10(6)
		timeUnit = ""
	} else if elapsedTime > math.Pow10(3) {
		elapsedTime = elapsedTime / math.Pow10(3)
		timeUnit = "m"

	}

	fmt.Printf("Compiled in %f %ss\n", elapsedTime, timeUnit)

	os.Chdir(prevDir)

	return AssembledProgram{

		ProgramBytes:     programBytes,
		ProgramStructure: programData,
		ProgramJson:      string(jsonBytes),
	}, nil

}
