package compiler

import (
	"fmt"
	"sccreeper/govm/pkg/constants"
	"strings"
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

	//Split code into array based on line breaks

	program_list := strings.Split(output_path, "\n")

	var program_instructions []instruction
	var definitions []definition
	var jump_blocks []jump_block

	in_jump_block := false
	in_element := false

	var program_statements [][]string

	//Begin parsing of statements

	for _, statement := range program_list {

		//Ignore if comment

		if statement[0:1] == "//" {
			continue
		}

		line := []rune(statement)
		var statement_items []string

		current_statement := ""

		//Loop to split the statement into individual elements (instructions, registers, data etc.)
		for _, char := range line {

			in_element = true

			if char == ' ' {

				statement_items = append(statement_items, current_statement)

			} else {
				current_statement += string(char)
				continue
			}

		}
		//Add the final item to the end of the statement
		statement_items = append(statement_items, current_statement)

		//Add the semi-parsed statement to the statement splice
		program_statements = append(program_statements, statement_items)

		in_element = false

	}

	//Debug, print statements to console

	for _, e := range program_statements {

		fmt.Println("Statement %v", e)

	}

}
