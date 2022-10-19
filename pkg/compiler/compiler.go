package compiler

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sccreeper/govm/pkg/constants"
	"sccreeper/govm/pkg/util"
	"strconv"
	"strings"
	"time"
)

// Compiler backend

//Config for CLI

type CompilerConfig struct {
	OutputPath string

	OutputJSON bool
	JSONPath   string
}

// Types for statements
type instruction struct {
	SingleData bool     `json:"single_data"`
	Data       []string `json:"data"`

	Instruction uint32 `json:"instruction"`
}

type definition struct {
	Name string            `json:"name"`
	Data []byte            `json:"data"`
	Type constants.DefType `json:"type"`
}

type interrupt_subscription struct {
	Interrupt     constants.Interrupt `json:"interrupt"`
	JumpBlockName string              `json:"jump_block_name"`
}

type jump_block struct {
	Name         string        `json:"name"`
	Instructions []instruction `json:"instructions"`
}

// Struct for holding program data
type program_structure struct {
	AllNames []string `json:"all_names"`

	JumpBlockNames         []string                 `json:"jump_block_names"`
	DefNames               []string                 `json:"definition_names"`
	InterruptSubscriptions []interrupt_subscription `json:"interrupt_subscriptions"`

	ProgramInstructions []instruction         `json:"program_instructions"`
	Definitions         []definition          `json:"definitions"`
	JumpBlocks          map[string]jump_block `json:"jump_blocks"`
}

var all_names []string

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

	if util.ContainsString(all_names, s) {
		err = fmt.Errorf("%s collides with %s", s, s)
	}

	return err
}

// Compile method
func Compile(code_string string, config CompilerConfig) error {

	start_time := time.Now().UnixMicro()

	//Remove empty lines

	code_string = strings.ReplaceAll(code_string, "\n\n", "\n")

	//Split code into array based on line breaks

	program_list := strings.Split(code_string, "\n")

	log.Printf("Code size: %d byte(s)", len(code_string))
	log.Printf("Code lines: %d line(s)", len(program_list))

	//---------------------------
	//Begin parsing of statements
	//---------------------------

	program_statements := make([][]string, 0)

	in_element := false
	_ = in_element
	in_string := false

	for index, statement := range program_list {

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

	log.Println("Finished first stage of parsing...")

	//Debug, print statements to console

	for _, e := range program_statements {

		log.Printf("Statement %s\n", e)
		//log.Printf("Statement data %s\n", e[0:])

	}

	//------------------------
	// Begin data construction
	//------------------------

	//Make program data struct

	var program_data = program_structure{
		JumpBlocks: make(map[string]jump_block),
	}

	//Arrays for names
	//var jump_block_names []string
	//var def_names []string
	//var interrupt_subscriptions []interrupt_subscription

	//Instructions not contained in any jump blocks
	//var program_instructions []instruction

	//var definitions []definition
	//var jump_blocks = make(map[string]jump_block)

	var current_jump_block_instructions []instruction
	jump_block_name := ""
	in_jump_block := false

	for index, e := range program_statements {

		log.Printf("Parsing statement %d", index)

		if len(e) == 0 {
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
				data_array = []byte(e[2])

			case constants.IntType:
				i, err := strconv.ParseInt(e[2], 10, 32)
				util.Check(err)

				buffer := new(bytes.Buffer)
				binary.Write(buffer, binary.LittleEndian, i)

				data_array = []byte(buffer.Bytes())
			}

			program_data.Definitions = append(program_data.Definitions,
				definition{
					Name: e[1],
					Data: data_array,
					Type: def_type,
				},
			)

		} else if e[0] == "sub" { //Interrupt subscription

			//Error checking

			if _, exists := constants.InterruptInts[e[1]]; !exists || constants.InterruptInts[e[1]] < constants.IntMouseMove {
				return fmt.Errorf("unrecognized interrupt %s", e[1])
			}

			if !util.ContainsString(program_data.JumpBlockNames, e[2]) {
				return fmt.Errorf("unreconized jump %s", e[2])
			}

			program_data.InterruptSubscriptions = append(
				program_data.InterruptSubscriptions,

				interrupt_subscription{
					Interrupt:     constants.Interrupt(constants.InterruptInts[e[1]]),
					JumpBlockName: e[2],
				},
			)

		} else if e[0] == "end" { //Reaching end of jump block
			if !in_jump_block {
				return errors.New("unexpected end statement")
			}

			program_data.JumpBlocks[jump_block_name] = jump_block{

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
				return errors.New("cannot nest jump blocks")
			}
			if len(e[0]) == 1 {
				return errors.New("jumpblock names must have minimum length of one")
			}
			//Check if name of jump block isn't shared by registers or instructions
			name_collision(e[0][1:])

			jump_block_name = e[0][1:]
			program_data.AllNames = append(program_data.AllNames, e[0][1:])

			in_jump_block = true
			program_data.JumpBlockNames = append(program_data.JumpBlockNames, e[0][1:])
			all_names = append(all_names, e[0][1:])

			continue

		} else {

			//Parse for other statements

			//Check if statement exists in instructions
			if _, exists := constants.InstructionInts[e[0]]; !exists {
				return fmt.Errorf("instruction %s does not exist", e[0])
			}

			//If does exist, continue

			single_data := false

			if len(e[1:]) == 1 {
				single_data = true
			}

			instruction_to_be_added := instruction{
				SingleData:  single_data,
				Data:        e[1:],
				Instruction: constants.InstructionInts[e[0]],
			}

			if in_jump_block {

				current_jump_block_instructions = append(current_jump_block_instructions, instruction_to_be_added)

			} else {
				program_data.ProgramInstructions = append(program_data.ProgramInstructions, instruction_to_be_added)
			}

		}
	}

	//-------------------------
	// Begin bytecode generation
	//--------------------------

	//Generate addresses for definitions (data block)

	//Defer finish compile time

	// -----------------
	// Output JSON
	// ----------------

	if config.OutputJSON {

		json_bytes, err := json.MarshalIndent(program_data, "", "\t")

		util.CheckError(err)

		err = os.WriteFile(config.JSONPath, json_bytes, 0644)

		util.CheckError(err)

		log.Printf("Outputted JSON structure to '%s'", config.JSONPath)

	}

	// -------------------
	// Output elapsed time
	// -------------------

	elasped_time := float64(time.Now().UnixMicro() - start_time)
	time_unit := "µ"

	if elasped_time > math.Pow10(6) {
		elasped_time = elasped_time / math.Pow10(6)
		time_unit = ""
	} else if elasped_time > math.Pow10(3) {
		elasped_time = elasped_time / math.Pow10(3)
		time_unit = "m"

	}

	log.Printf("Compiled in %f %ss", elasped_time, time_unit)

	return nil

}
