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
	InterruptName string              `json:"interrupt_name"`
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

//Constants

const (
	InstructionLength uint32 = 5 //Instruction length in bytes
	BlockAddrSize            = 4 * 4
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

	if util.ContainsString(all_names, s) {
		err = fmt.Errorf("%s collides with %s", s, s)
	}

	return err
}

//General purpose instruction for generating instruction bytecode

func generate_instruction_bytecode(i instruction, d_block_addr map[string]uint32, j_blk_addr map[string]uint32) []byte {

	//TODO: sign bit
	//TODO: add offset for "hardware reserved" space

	var instruction_bytes []byte

	instruction_bytes = append(instruction_bytes,
		uint8(i.Instruction),
	)

	//Evaluate instruction args

	var addresses []uint32

	for _, v := range i.Data {
		var addr uint32

		if v[0] == '@' {

			addr = j_blk_addr[v[1:]]

		} else if i.Instruction == uint32(constants.IJump) {

			addr = uint32(j_blk_addr[v])

		} else {
			addr = uint32(constants.InterruptInts[v])
		}

		addresses = append(addresses, addr)
	}

	//Add args to byte array

	var data_array []byte

	if i.SingleData {

		data_array = make([]byte, 4)

		binary.LittleEndian.PutUint32(data_array[:], addresses[0])

	} else {
		data_array = make([]byte, 4)

		binary.LittleEndian.PutUint16(data_array[:], uint16(addresses[0]))
		binary.LittleEndian.PutUint16(data_array[2:], uint16(addresses[1]))
	}

	instruction_bytes = append(instruction_bytes, data_array...)

	return instruction_bytes

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
					InterruptName: e[1],
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

	log.Println("Starting bytecode generation")
	log.Println("Starting data block generation")

	//Start generating the data block.

	var byte_index uint32 = BlockAddrSize

	data_block_addresses := make(map[string]uint32)
	data_block_bytes := []byte{}

	data_start_index := byte_index

	for _, v := range program_data.Definitions {

		data_block_addresses[v.Name] = byte_index

		data_array := make([]byte, 4)

		binary.LittleEndian.PutUint32(data_array[:], uint32(len(v.Data)))

		data_block_bytes = append(data_block_bytes, data_array...)
		data_block_bytes = append(data_block_bytes, v.Data...)

		byte_index += uint32(4 + len(v.Data))
		//                   4 - Length of uint32 which indicates length
		//                     - Length of v.Data (which is already in bytes)

	}

	//Generate jump block bytecode and instructions

	log.Println("Starting jump block generation")

	jump_block_addresses := make(map[string]uint32)
	jump_block_bytes := []byte{}

	var instruction_byte_array []byte

	jmp_block_start_index := byte_index

	for _, jmp_blk := range program_data.JumpBlocks {

		jump_block_addresses[jmp_blk.Name] = byte_index

		for _, jmp_i := range jmp_blk.Instructions {

			instruction_byte_array = append(
				instruction_byte_array,
				generate_instruction_bytecode(
					jmp_i,
					data_block_addresses,
					jump_block_addresses,
				)...,
			)

		}

		//Add null terminator to jmp
		instruction_byte_array = append(instruction_byte_array, []byte{0, 0, 0, 0}...)

		//Add jump block to bytes
		jump_block_bytes = append(jump_block_bytes, instruction_byte_array...)

		byte_index += uint32((len(instruction_byte_array) - 1) + int(InstructionLength))

	}

	//Generate interrupt table

	interrupt_table_start_index := byte_index

	var interrupt_bytes []byte

	interrupts := make(map[string]interrupt_subscription)

	//Generate interrupts map

	for _, v := range program_data.InterruptSubscriptions {

		interrupts[v.InterruptName] = v

	}

	var int_address uint32

	for k, v := range constants.SubscribableInterrupts {

		if _, exists := interrupts[k]; exists {

			int_address = jump_block_addresses[interrupts[k].JumpBlockName]

		} else {
			int_address = uint32(0)
		}

		data_array := make([]byte, 8)

		binary.LittleEndian.PutUint32(data_array[:], uint32(v))
		binary.LittleEndian.PutUint32(data_array[4:], uint32(int_address))

		interrupt_bytes = append(interrupt_bytes, data_array...)

		byte_index += uint32(len(data_array))

	}

	//---------------

	//Finally, generate other instructions

	log.Println("Starting other instruction generation")

	var other_instruction_bytes []byte

	instruction_start_index := byte_index

	for _, v := range program_data.ProgramInstructions {

		other_instruction_bytes = append(
			other_instruction_bytes,
			generate_instruction_bytecode(
				v,
				data_block_addresses,
				jump_block_addresses,
			)...,
		)

		byte_index += InstructionLength

	}

	//--------------------------
	//Build final program binary
	//--------------------------

	var final_byte_array []byte

	//Start with block indexes

	block_index_array := make([]byte, 16)

	binary.LittleEndian.PutUint32(block_index_array[:], data_start_index)
	binary.LittleEndian.PutUint32(block_index_array[4:], jmp_block_start_index)
	binary.LittleEndian.PutUint32(block_index_array[8:], interrupt_table_start_index)
	binary.LittleEndian.PutUint32(block_index_array[12:], instruction_start_index)

	final_byte_array = append(final_byte_array, block_index_array...)

	//Add data, jumps, interrupts, and program

	final_byte_array = append(final_byte_array, data_block_bytes...)
	final_byte_array = append(final_byte_array, jump_block_bytes...)
	final_byte_array = append(final_byte_array, interrupt_bytes...)
	final_byte_array = append(final_byte_array, other_instruction_bytes...)

	//Write to file

	os.WriteFile(config.OutputPath, final_byte_array, 0666)
	//Output start indexes

	log.Printf("Data start index: %d", data_start_index)
	log.Printf("Jump start index: %d", jmp_block_start_index)
	log.Printf("Interrupt table start index: %d", interrupt_table_start_index)
	log.Printf("Program start index: %d", instruction_start_index)
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
