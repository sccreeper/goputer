package compiler

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"
	"strconv"
	"strings"
)

var all_names []string

//Name collision function

func name_collision(s string) error {

	var err error = nil

	if _, exists := constants.InstructionInts[s]; exists {
		err = fmt.Errorf("name %s shares name with instruction", s)
	}
	if _, exists := constants.RegisterInts[s]; exists {
		err = fmt.Errorf("name %s shares name with register", s)
	}
	if _, exists := constants.InterruptInts[s]; exists {
		err = fmt.Errorf("name %s shares name with interrupt", s)
	}

	if util.SliceContains(all_names, s) {
		err = fmt.Errorf("%s collides with %s", s, s)
	}

	return err
}

// Takes a string and returns a program structure
//
// code_string - code string
//
// verbose - wether to output console when parsing
//
// imported - are we importing this file from another file
func parse(code_string string, verbose bool, imported bool) (ProgramStructure, error) {

	//Remove empty lines

	code_string = strings.ReplaceAll(code_string, "\n\n", "\n")

	//Split code into array based on line breaks

	program_list := strings.Split(code_string, "\n")

	if verbose {

		log.Printf("Code size: %d byte(s)", len(code_string))
		log.Printf("Code lines: %d line(s)", len(program_list))
	}

	//---------------------------
	//Begin parsing of statements
	//---------------------------

	program_statements := make([][]string, 0)

	in_element := false
	_ = in_element
	in_string := false

	for index, statement := range program_list {

		if statement == "" {
			continue
		}

		in_string = false

		//Ignore if comment

		if statement[:2] == "//" {
			program_statements = append(program_statements, nil)
			continue
		}

		line := statement
		//statement_items := make([]string, 0)
		current_statement := ""

		//Remove trailing whitespace

		line = strings.Trim(line, " ")

		//Loop to split the statement into individual elements (instructions, registers, data etc.)
		for _, char := range line {

			in_element = true

			current_statement += string(char)

			if (char == ' ' && !in_string) || (in_string && char == '"') {

				if len(program_statements)-1 < index || index == 0 {

					program_statements = append(program_statements, make([]string, 0))
				}

				if char == '"' {
					in_string = false
				}

				program_statements[index] = append(program_statements[index], strings.TrimSpace(current_statement))

				current_statement = ""

			}

			if char == '"' {
				in_string = true
			}

			//Add the semi-parsed statement to the statement splice
			//program_statements = append(program_statements, statement_items)

		}

		if len(program_statements)-1 < index || index == 0 {

			program_statements = append(program_statements, make([]string, 0))
		}

		program_statements[index] = append(program_statements[index], strings.TrimSpace(current_statement))
	}

	if verbose {
		log.Println("Finished first stage of parsing...")

		//Debug, print statements to console

		for _, e := range program_statements {

			log.Printf("Statement %s\n", e)
			//log.Printf("Statement data %s\n", e[0:])

		}

	}

	var program_data = ProgramStructure{
		InstructionBlocks: make(map[string]CodeBlock),
	}

	//--------------------------
	//Evaluate import statements
	//--------------------------

	for _, e := range program_statements {

		if len(e) < 2 {
			continue
		} else if e[0] == "import" {

			//Read other file
			f_name := strings.Trim(e[1], "\"")

			imported_file, err := os.ReadFile(f_name)
			util.CheckError(err)

			imported_program_structure, err := parse(string(imported_file), false, true)
			util.CheckError(err)

			program_data, err = combine(program_data, imported_program_structure)
			util.CheckError(err)

		}

	}

	//------------------------
	// Begin data construction
	//------------------------

	//Make program data struct

	var current_jump_block_instructions []Instruction
	jump_block_name := ""
	in_jump_block := false

	for index, e := range program_statements {

		if verbose {
			log.Printf("Parsing statement %d", index)
		}

		if len(e) == 0 {
			continue
		}

		if e[0] == "import" {
			continue
		}

		// Parse for special purpose statements

		if e[0] == "def" { //Constant definition
			name_collision(e[1])

			program_data.DefNames = append(program_data.DefNames, e[1])
			program_data.AllNames = append(program_data.AllNames, e[1])

			// Parse definition data, decide wether is int string, float, etc.

			var def_type constants.DefType = 0
			data_array := make([]byte, 4)

			//Is float
			if strings.Contains(e[2], ".") && !(e[2][0] == '"') {
				def_type = constants.FloatType
			} else if e[2][0] == '-' { //Signed int
				def_type = constants.IntType
			} else if e[2][0] == '"' {
				def_type = constants.StringType
			} else if len(e[2]) > 2 && e[2][0:2] == "0x" || e[2][0] == '@' {
				def_type = constants.BytesType
			} else {
				def_type = constants.UintType
			}

			//Convert definition data to byte array
			switch def_type {
			case constants.FloatType:
				i, err := strconv.ParseFloat(e[2], 32)
				util.Check(err)
				binary.LittleEndian.PutUint32(data_array[:], math.Float32bits(float32(i)))

			case constants.UintType:
				i, err := strconv.ParseUint(e[2], 10, 32)
				util.Check(err)
				binary.LittleEndian.PutUint32(data_array[:], uint32(i))

			case constants.StringType:
				//Remove speech marks

				e[2] = strings.Trim(e[2], "\"")

				data_array = []byte(e[2])

			case constants.IntType:
				i, err := strconv.ParseInt(e[2], 10, 32)
				util.Check(err)

				buffer := new(bytes.Buffer)
				binary.Write(buffer, binary.LittleEndian, i)

				data_array = []byte(buffer.Bytes())
			case constants.BytesType:
				if e[2][0] == '@' {

					b, err := os.ReadFile(e[2][1:])
					util.CheckError(err)

					data_array = b

				} else {
					data_array, _ = hex.DecodeString(e[2][2:])
				}

			}

			program_data.Definitions = append(program_data.Definitions,
				Definition{
					Name: e[1],
					Data: data_array,
					Type: def_type,
				},
			)

		} else if e[0] == "sub" { //Interrupt subscription

			//Error checking

			if _, exists := constants.InterruptInts[e[1]]; !exists || constants.InterruptInts[e[1]] < constants.IntMouseMove {
				return program_data, fmt.Errorf("unrecognized interrupt %s", e[1])
			}

			if !util.SliceContains(program_data.InstructionBlockNames, e[2]) {
				return program_data, fmt.Errorf("unreconized jump %s", e[2])
			}

			program_data.InterruptSubscriptions = append(
				program_data.InterruptSubscriptions,

				InterruptSubscription{
					InterruptName: e[1],
					Interrupt:     constants.Interrupt(constants.InterruptInts[e[1]]),
					JumpBlockName: e[2],
				},
			)

		} else if e[0] == "end" { //Reaching end of jump block
			if !in_jump_block {
				return program_data, errors.New("unexpected end statement")
			}

			program_data.InstructionBlocks[jump_block_name] = CodeBlock{

				Name:         jump_block_name,
				Instructions: current_jump_block_instructions,
			}

			in_jump_block = false
			jump_block_name = ""
			current_jump_block_instructions = nil

			continue

		} else if e[0][0] == ':' {
			//Errors
			if in_jump_block {
				return program_data, errors.New("cannot nest jump blocks")
			}
			if len(e[0]) == 1 {
				return program_data, errors.New("jumpblock names must have minimum length of one")
			}
			//Check if name of jump block isn't shared by registers or instructions
			name_collision(e[0][1:])

			jump_block_name = e[0][1:]
			program_data.AllNames = append(program_data.AllNames, e[0][1:])

			in_jump_block = true
			program_data.InstructionBlockNames = append(program_data.InstructionBlockNames, e[0][1:])
			all_names = append(all_names, e[0][1:])

			continue

		} else {

			//Parse for other statements

			//Check if statement exists in instructions
			if _, exists := constants.InstructionInts[e[0]]; !exists {
				return program_data, fmt.Errorf("instruction %s does not exist", e[0])
			}

			//If does exist, continue

			single_data := false

			if len(e[1:]) == 1 {
				single_data = true
			}

			instruction_to_be_added := Instruction{
				SingleData:  single_data,
				Data:        e[1:],
				Instruction: constants.InstructionInts[e[0]],
			}

			if in_jump_block {

				current_jump_block_instructions = append(current_jump_block_instructions, instruction_to_be_added)

			} else if !imported {
				program_data.ProgramInstructions = append(program_data.ProgramInstructions, instruction_to_be_added)
			}

		}
	}

	return program_data, nil
}

func combine(s0 ProgramStructure, s1 ProgramStructure) (ProgramStructure, error) {

	var combined ProgramStructure

	//Merge splices & check for name conflicts

	for _, v := range s0.AllNames {

		if util.SliceContains(s1.AllNames, v) {
			return ProgramStructure{}, fmt.Errorf("name conflict involving '%s'", v)
		}

	}

	combined.AllNames = append(s0.AllNames, s1.AllNames...)

	combined.InstructionBlockNames = append(s0.InstructionBlockNames, s1.InstructionBlockNames...)
	combined.DefNames = append(s0.DefNames, s1.DefNames...)

	combined.Definitions = append(s0.Definitions, s1.Definitions...)

	//Combine instruction blocks

	combined.InstructionBlocks = s0.InstructionBlocks

	for k, v := range s1.InstructionBlocks {

		combined.InstructionBlocks[k] = v

	}

	return combined, nil

}
