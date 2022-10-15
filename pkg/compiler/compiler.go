package compiler

import (
	"log"
	"sccreeper/govm/pkg/constants"
	"strings"
	"time"
)

// Compiler backend

type instruction struct {
	single_data bool
	d0          uint32
	d1          uint32

	instruction uint32
}

type definition struct {
	name string
	d0   []byte
}

type jump_block struct {
	name        string
	instruction constants.Instruction
}

var jump_block_names []string
var def_names []string

func Compile(code_string string, output_path string) {

	start_time := time.Now().UnixMicro()

	//Remove empty lines

	code_string = strings.Replace(code_string, "\n\n", "\n", -1)

	//Split code into array based on line breaks

	program_list := strings.Split(code_string, "\n")

	log.Printf("Code size: %d byte(s)", len(code_string))
	log.Printf("Code lines: %d line(s)", len(program_list))

	// var program_instructions []instruction
	// var definitions []definition
	// var jump_blocks []jump_block

	// in_jump_block := false

	in_element := false
	_ = in_element

	program_statements := make([][]string, 0)

	//Begin parsing of statements

	for index, statement := range program_list {

		//Ignore if comment

		if statement[:2] == "//" {
			program_statements = append(program_statements, nil)
			continue
		}

		line := []rune(statement)
		//statement_items := make([]string, 0)
		current_statement := ""

		//Loop to split the statement into individual elements (instructions, registers, data etc.)
		for _, char := range line {

			in_element = true

			if char == ' ' {

				if len(program_statements)-1 < index || index == 0 {

					program_statements = append(program_statements, make([]string, 0))
				}

				program_statements[index] = append(program_statements[index], current_statement)

				current_statement = ""

			} else {
				current_statement += string(char)
				continue
			}
		}
		//Add the final item to the end of the statement

		if len(program_statements)-1 < index || index == 0 {

			program_statements = append(program_statements, make([]string, 0))
		}

		program_statements[index] = append(program_statements[index], current_statement)

		//Add the semi-parsed statement to the statement splice
		//program_statements = append(program_statements, statement_items)

	}

	log.Println("Finished first stage of parsing...")

	//Debug, print statements to console

	for _, e := range program_statements {

		log.Printf("Statement %s\n", e)
		//log.Printf("Statement data %s\n", e[0:])

	}

	// Print elapsed time
	log.Printf("Compiled in %d microseconds(s)", time.Now().UnixMicro()-start_time)

}
