package compiler

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"os"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"
	"strconv"
	"strings"
)

type Parser struct {
	CodeString string
	CodeLines  []string

	LineIndex   int
	CurrentLine string

	AllNames []string

	ProgramStatements [][]string
	ProgramStructure  ProgramStructure

	FileName string
	Imported bool

	Verbose bool

	ErrorHandler func(error_type ErrorType, error_text string)
}

// Takes a string and returns a program structure
//
// code_string - code string
//
// verbose - wether to output console when parsing
//
// imported - are we importing this file from another file

func (p *Parser) Parse() (ProgramStructure, error) {

	//Remove empty lines

	//p.CodeString = strings.ReplaceAll(p.CodeString, "\n\n", "\n")

	//Split code into array based on line breaks

	p.CodeLines = strings.Split(p.CodeString, "\n")

	if p.Verbose {

		log.Printf("Code size: %d byte(s)", len(p.CodeLines))
		log.Printf("Code lines: %d line(s)", len(p.CodeLines))
	}

	//---------------------------
	//Begin parsing of statements
	//---------------------------

	p.ProgramStatements = make([][]string, 0)

	in_element := false
	_ = in_element
	in_string := false

	for index, statement := range p.CodeLines {

		in_string = false

		if statement == "" {
			p.ProgramStatements = append(p.ProgramStatements, []string{})
			continue
		} else if statement == "\n" {
			p.ProgramStatements = append(p.ProgramStatements, []string{})
			continue
		} else if statement[:2] == "//" { //Ignore if comment
			p.ProgramStatements = append(p.ProgramStatements, []string{statement})
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

				if len(p.ProgramStatements)-1 < index || index == 0 {

					p.ProgramStatements = append(p.ProgramStatements, make([]string, 0))
				}

				if char == '"' {
					in_string = false
				}

				p.ProgramStatements[index] = append(p.ProgramStatements[index], strings.TrimSpace(current_statement))

				current_statement = ""

			}

			if char == '"' {
				in_string = true
			}

			//Add the semi-parsed statement to the statement splice
			//program_statements = append(program_statements, statement_items)

		}

		if len(p.ProgramStatements)-1 < index || index == 0 {

			p.ProgramStatements = append(p.ProgramStatements, make([]string, 0))
		}

		p.ProgramStatements[index] = append(p.ProgramStatements[index], strings.TrimSpace(current_statement))
	}

	if p.Verbose {
		log.Println("Finished first stage of parsing...")

		//Debug, print statements to console

		for _, e := range p.ProgramStatements {

			log.Printf("Statement %s\n", e)
			log.Printf("Statement length %d\n", len(e))

		}

	}

	p.ProgramStructure = ProgramStructure{
		InstructionBlocks: make(map[string]CodeBlock),
		ImportedFiles:     []string{p.FileName},
	}

	//--------------------------
	//Evaluate import statements
	//--------------------------

	for _, e := range p.ProgramStatements {

		if len(e) < 2 {
			continue
		} else if e[0] == "import" {

			//Read other file
			f_name := strings.Trim(e[1], "\"")

			if util.SliceContains(p.ProgramStructure.ImportedFiles, f_name) {
				//Already imported
				continue
			}

			imported_file, err := os.ReadFile(f_name)

			if err == os.ErrNotExist {
				p.parsing_error(ErrImport, FileNotFound)
			} else {
				util.CheckError(err)
			}

			import_parser := Parser{
				CodeString:   string(imported_file),
				FileName:     p.FileName,
				Imported:     true,
				Verbose:      false,
				ErrorHandler: p.ErrorHandler,
			}

			imported_program_structure, err := import_parser.Parse()
			util.CheckError(err)

			p.ProgramStructure, err = p.combine(imported_program_structure)
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

	for index, e := range p.ProgramStatements {

		p.LineIndex = index
		p.CurrentLine = p.CodeLines[index]

		if p.Verbose {
			log.Printf("Parsing statement %d", index)
		}

		if len(e) == 0 {
			continue
		} else if len(e) == 1 && e[0] == "" {
			continue
		} else if e[0] == "import" {
			continue
		} else if e[0][:2] == "//" {
			continue
		}

		// Parse for special purpose statements

		if e[0] == "def" { //Constant definition
			p.name_collision(e[1])

			p.ProgramStructure.DefNames = append(p.ProgramStructure.DefNames, e[1])
			p.ProgramStructure.AllNames = append(p.ProgramStructure.AllNames, e[1])

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
				e[2] = strings.Replace(e[2], `\n`, "\n", -1)

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

			p.ProgramStructure.Definitions = append(p.ProgramStructure.Definitions,
				Definition{
					Name: e[1],
					Data: data_array,
					Type: def_type,
				},
			)

			continue

		} else if e[0] == "sub" { //Interrupt subscription

			//Error checking

			if _, exists := constants.InterruptInts[e[1]]; !exists || constants.InterruptInts[e[1]] < constants.IntMouseMove {
				p.parsing_error(ErrSymbol, SymbolDoesNotExist)
			}

			if !util.SliceContains(p.ProgramStructure.InstructionBlockNames, e[2]) {
				p.parsing_error(ErrSymbol, ErrorType(fmt.Sprintf("unrecognized jump %s", e[2])))
			}

			p.ProgramStructure.InterruptSubscriptions = append(
				p.ProgramStructure.InterruptSubscriptions,

				InterruptSubscription{
					InterruptName: e[1],
					Interrupt:     constants.Interrupt(constants.InterruptInts[e[1]]),
					JumpBlockName: e[2],
				},
			)

			continue

		} else if e[0] == "end" { //Reaching end of jump block
			if !in_jump_block {
				p.parsing_error(ErrSyntax, UnexpectedEndStatement)
			}

			p.ProgramStructure.InstructionBlocks[jump_block_name] = CodeBlock{

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
				p.parsing_error(ErrSyntax, NestingError)
			}
			if len(e[0]) == 1 {
				p.parsing_error(ErrSyntax, MinimumNameLength)
			}
			//Check if name of jump block isn't shared by registers or instructions
			p.name_collision(e[0][1:])

			jump_block_name = e[0][1:]
			p.ProgramStructure.AllNames = append(p.ProgramStructure.AllNames, e[0][1:])

			in_jump_block = true
			p.ProgramStructure.InstructionBlockNames = append(p.ProgramStructure.InstructionBlockNames, e[0][1:])

			continue

		} else {

			//Parse for other statements

			//Check if statement exists in instructions
			if _, exists := constants.InstructionInts[e[0]]; !exists {
				p.parsing_error(ErrDoesNotExist, InstructionDoesNotExist)
			}

			//Check if args are valid

			for _, arg := range e[1:] {

				if e[0] == "int" {
					if _, exists := constants.InterruptInts[arg]; !exists {
						p.parsing_error(ErrSymbol, InvalidArgument)
					}

				} else if arg[0] == '@' && (e[0] == "lda" || e[0] == "sta") {

					var exists bool = false

					for _, v := range p.ProgramStructure.DefNames {

						if arg[1:] == v {
							exists = true
							break
						}
					}

					if !exists {
						p.parsing_error(ErrDoesNotExist, ErrorType(fmt.Sprintf("definition '%s' does not exist", e[1][1:])))
					}
				} else if e[0] == "jmp" || e[0] == "cndjmp" || e[0] == "call" || e[0] == "cndcall" {

					var exists bool = false

					for _, v := range p.ProgramStructure.InstructionBlockNames {

						if e[1] == v {
							exists = true
							break
						}

					}

					if !exists {
						p.parsing_error(ErrDoesNotExist, ErrorType(fmt.Sprintf("unknown instruction block '%s'", arg)))
					}

				} else {

					if _, exists := constants.RegisterInts[arg]; !exists {

						p.parsing_error(ErrDoesNotExist, ErrorType(fmt.Sprintf("unknown register '%s'", arg)))

					}

				}

			}

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

		} else if !p.Imported {
			p.ProgramStructure.ProgramInstructions = append(p.ProgramStructure.ProgramInstructions, instruction_to_be_added)
		}

	}

	return p.ProgramStructure, nil
}

// Method for combining program structures.
//
// Used for imports.
func (p *Parser) combine(s1 ProgramStructure) (ProgramStructure, error) {

	var combined ProgramStructure

	//Merge splices & check for name conflicts

	for _, v := range p.ProgramStructure.AllNames {

		if util.SliceContains(s1.AllNames, v) {
			return ProgramStructure{}, ErrSymbol
		}

	}

	if len(s1.AllNames) > 0 {
		combined.AllNames = append(p.ProgramStructure.AllNames, s1.AllNames...)
	} else {
		combined.AllNames = p.ProgramStructure.AllNames
	}
	if len(s1.InstructionBlockNames) > 0 {
		combined.InstructionBlockNames = append(p.ProgramStructure.InstructionBlockNames, s1.InstructionBlockNames...)
	} else {
		combined.InstructionBlockNames = p.ProgramStructure.InstructionBlockNames
	}
	if len(s1.AllNames) > 0 {
		combined.DefNames = append(p.ProgramStructure.DefNames, s1.DefNames...)
	} else {
		combined.DefNames = p.ProgramStructure.DefNames
	}

	if len(s1.Definitions) > 0 {
		combined.Definitions = append(p.ProgramStructure.Definitions, s1.Definitions...)
	} else {
		combined.Definitions = p.ProgramStructure.Definitions
	}

	//Combine instruction blocks

	combined.InstructionBlocks = p.ProgramStructure.InstructionBlocks

	if len(s1.InstructionBlocks) > 0 {
		for k, v := range s1.InstructionBlocks {

			combined.InstructionBlocks[k] = v

		}
	}

	//Check for circular imports

	if util.SliceContains(s1.ImportedFiles, p.FileName) {

		p.parsing_error(ErrImport, CircularImport)

	}

	//Combine imports

	for _, v := range s1.ImportedFiles {

		if !util.SliceContains(p.ProgramStructure.ImportedFiles, v) {
			combined.ImportedFiles = append(combined.ImportedFiles, v)
		}

	}

	combined.ImportedFiles = append(p.ProgramStructure.ImportedFiles, s1.ImportedFiles...)

	return combined, nil

}

//Name collision function

func (p *Parser) name_collision(s string) string {

	var err string = ""

	if _, exists := constants.InstructionInts[s]; exists {
		err = fmt.Sprintf("name %s shares name with instruction", s)
	}
	if _, exists := constants.RegisterInts[s]; exists {
		err = fmt.Sprintf("name %s shares name with register", s)
	}
	if _, exists := constants.InterruptInts[s]; exists {
		err = fmt.Sprintf("name %s shares name with interrupt", s)
	}

	if util.SliceContains(p.ProgramStructure.AllNames, s) {
		err = fmt.Sprintf("%s collides with %s", s, s)
	}

	if err != "" {
		p.parsing_error(ErrSymbol, ErrorType(err))
		return ""
	} else {
		return err
	}
}
