package compiler

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

// Regex expressions

var compTimeStatementRegex *regexp.Regexp
var commentStatementRegex *regexp.Regexp
var doubleQuoteStringValueRegex *regexp.Regexp
var intValueRegex *regexp.Regexp
var floatValueRegex *regexp.Regexp
var hexValueRegex *regexp.Regexp
var specialValueRegex *regexp.Regexp
var nameValueRegex *regexp.Regexp

func init() {

	// Regex used for validating expressions

	// Match statements that follow the form "#label value"
	compTimeStatementRegex = regexp.MustCompile(`^#([a-zA-Z]+) +(.+)$`)

	// Match comments that are at the beginning of a line
	commentStatementRegex = regexp.MustCompile(`^\/\/\s*(.*)`)

	// String value
	doubleQuoteStringValueRegex = regexp.MustCompile(`^"(?:\\.|[^"\\])*"$`)

	// Match integer values
	intValueRegex = regexp.MustCompile(`^(-?[0-9]+)$`)

	// Match floating point values
	floatValueRegex = regexp.MustCompile(`^(-?[0-9]+\.[0-9]+)$`)

	// Match hex values
	hexValueRegex = regexp.MustCompile(`^0x([a-zA-Z0-9]+)$`)

	// Match values in the form hello:world
	specialValueRegex = regexp.MustCompile(`^(.+):(.+)$`)

	// Match name value pairs e.g.
	// this_number 45
	// this string "isn't valid"

	nameValueRegex = regexp.MustCompile(`^([a-zA-Z0-9_\-]+)\s+(.+)$`)

}

// Parser struct
// Can parse files with the Parse() method.
type Parser struct {
	CodeString string   //Code of the entrypoint file.
	CodeLines  []string //Individual lines of code, do not set unless you know what you're doing.

	LineIndex   int    //Line index of the parser, do not set unless you know what you're doing.
	CurrentLine string //Current line of the parser, do not set unless you know what you're doing.

	AllNames []string //All the names in the program, do not set unless you know what you're doing.

	ProgramStatements [][]string       //Program lines, do not set unless you know what you're doing.
	ProgramStructure  ProgramStructure //The returned program structure, do not set unless you know what you're doing.

	FileName     string //Name of the entrypoint file (all imports are done relative to this.)
	Imported     bool   //Has this file been imported by another file (do not set this value unless you know what you are doing)
	ImportedFrom string //Where this file has been imported from (do not set this value unless you know what you are doing)

	Verbose bool //Should there be output when parsing the file.

	ErrorHandler func(error_type ErrorType, error_text string) //The error handler method, required.
	FileReader   func(path string) ([]byte, error)             //The file reader method, required.
}

// Parser method
//
// Needs to be called using a instantiated Parser struct.
//
// Returns a ProgramStructure struct.
func (p *Parser) Parse() (ProgramStructure, error) {

	//Split code into array based on line breaks

	p.CodeLines = strings.Split(p.CodeString, "\n")

	if p.Verbose {

		log.Printf("Code size: %d byte(s)", len(p.CodeLines))
		log.Printf("Code lines: %d line(s)", len(p.CodeLines))
	}

	// Begin parsing of statements

	p.ProgramStructure = ProgramStructure{
		ImportedFiles:          make([]string, 0),
		InterruptSubscriptions: make(map[string]InterruptSubscription),
		AllNames:               make([]string, 0),
		LabelNames:             make([]string, 0),
		DefinitionNames:        make([]string, 0),
		ProgramLabels:          make(map[string]ProgramLabel),
		ProgramInstructions:    make([]Instruction, 0),
		Definitions:            make(map[string]Definition),
	}

	// Pre-fill definitions and label names

	for _, line := range p.CodeLines {
		if compTimeStatementRegex.MatchString(line) {
			statementType := compTimeStatementRegex.FindStringSubmatch(line)[1]
			statementValue := compTimeStatementRegex.FindStringSubmatch(line)[2]

			if statementType == "label" {
				p.ProgramStructure.LabelNames = append(p.ProgramStructure.LabelNames, statementValue)
			} else if statementType == "def" {
				p.ProgramStructure.DefinitionNames = append(
					p.ProgramStructure.DefinitionNames,
					nameValueRegex.FindStringSubmatch(statementValue)[1],
				)
			}
		}
	}

	// Begin data construction

	var instructionCount int = 0

	for index, line := range p.CodeLines {

		p.LineIndex = index
		p.CurrentLine = p.CodeLines[index]

		if p.Verbose {
			log.Printf("Parsing statement %d", index)
		}

		// Skip conditions

		if len(line) == 0 || commentStatementRegex.MatchString(line) || len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// Keep going

		if compTimeStatementRegex.MatchString(line) {

			statementType := compTimeStatementRegex.FindStringSubmatch(line)[1]
			statementValue := compTimeStatementRegex.FindStringSubmatch(line)[2]

			switch statementType {
			case "label":

				if errMessage, collision := p.nameCollision(statementValue); collision {
					p.parsingError(ErrSymbol, ErrorType(errMessage))
				} else {
					p.AllNames = append(p.AllNames, statementValue)

					p.ProgramStructure.ProgramLabels[statementValue] = ProgramLabel{
						Name:              statementValue,
						InstructionOffset: instructionCount,
					}
				}

			case "def":

				defName := nameValueRegex.FindStringSubmatch(statementValue)[1]
				defStringValue := nameValueRegex.FindStringSubmatch(statementValue)[2]

				if errMessage, collison := p.nameCollision(defName); collison {
					p.parsingError(ErrSymbol, ErrorType(errMessage))
				} else {
					p.AllNames = append(p.AllNames, defName)

					var defByteValue []byte
					var defType constants.DefType

					if doubleQuoteStringValueRegex.MatchString(defStringValue) {

						strValue, err := strconv.Unquote(defStringValue)
						if err != nil {
							p.parsingError(ErrSyntax, ErrorType(err.Error()))
						}

						defByteValue = []byte(strValue)
						defType = constants.StringType

					} else if intValueRegex.MatchString(defStringValue) {

						defByteValue = make([]byte, 4)

						x, err := strconv.Atoi(defStringValue)
						if err != nil {
							p.parsingError(ErrSyntax, ErrorType(err.Error()))
						}

						binary.LittleEndian.PutUint32(defByteValue[:], uint32(x))
						defType = constants.IntType

					} else if floatValueRegex.MatchString(defStringValue) {

						x, err := strconv.ParseFloat(defStringValue, 32)
						if err != nil {
							p.parsingError(ErrSyntax, ErrorType(err.Error()))
						}

						binary.LittleEndian.PutUint32(defByteValue[:], math.Float32bits(float32(x)))
						defType = constants.FloatType

					} else if hexValueRegex.MatchString(defStringValue) {

						// Get actual hex value

						hexValue := hexValueRegex.FindStringSubmatch(defStringValue)[1]

						var err error
						defByteValue, err = hex.DecodeString(hexValue)
						if err != nil {
							panic(err)
						}

						defType = constants.BytesType

					} else if specialValueRegex.MatchString(defStringValue) {

						specialType := specialValueRegex.FindStringSubmatch(defStringValue)[1]
						specialValue := specialValueRegex.FindStringSubmatch(defStringValue)[2]

						if specialType == "file" {

							// See if value is string

							if doubleQuoteStringValueRegex.MatchString(specialValue) {

								// Remove speechmarks

								specialValue = specialValue[1 : len(specialValue)-1]

								b, err := os.ReadFile(specialValue)
								if err != nil {
									panic(err)
								}

								defByteValue = b
								defType = constants.BytesType

							} else {
								p.parsingError(ErrSyntax, ErrorType("malformed string"))
							}

						} else if specialType == "region" {

							// Get region size

							regionSize, err := strconv.Atoi(specialValue)
							if err != nil {
								p.parsingError(ErrValue, InvalidValue)
							}

							if regionSize < 0 {
								p.parsingError(ErrValue, "region size cannot be negative")
							}

							defByteValue = make([]byte, regionSize, regionSize)
							defType = constants.BytesType

						} else {
							p.parsingError(ErrSyntax, ErrorType(fmt.Sprintf("unrecognised special definition type '%s'", specialType)))
						}

					} else {
						p.parsingError(ErrSyntax, InvalidValue)
					}

					p.ProgramStructure.Definitions[defName] = Definition{
						Name:       defName,
						StringData: defStringValue,
						ByteData:   defByteValue,
						Type:       defType,
					}
				}

			case "import":

				if !doubleQuoteStringValueRegex.MatchString(statementValue) {
					p.parsingError(ErrSyntax, "value should be string")
				}

				//Read other file

				fName := strings.Trim(statementValue, "\"")

				if slices.Contains(p.ProgramStructure.ImportedFiles, fName) {
					continue
				} else if fName == p.ImportedFrom {
					p.parsingError(ErrImport, CircularImport)
				}

				importedFile, err := p.FileReader(fName)

				if err != nil {
					p.parsingError(ErrFile, ErrorType(fmt.Sprintf("error reading file '%s'", fName)))
				}

				importParser := Parser{
					CodeString:   string(importedFile[:]),
					FileName:     fName,
					Imported:     true,
					ImportedFrom: p.FileName,
					Verbose:      false,
					ErrorHandler: p.ErrorHandler,
					FileReader:   p.FileReader,
				}

				p.ProgramStructure.ImportedFiles = append(p.ProgramStructure.ImportedFiles, fName)

				importedProgramStructure, err := importParser.Parse()
				if err != nil {
					p.parsingError(err, ErrorType(err.Error()))
				}

				if p.ImportedFrom == fName {
					p.parsingError(ErrImport, CircularImport)
				}

				p.ProgramStructure, err = p.combine(importedProgramStructure)
				if err != nil {
					p.parsingError(err, ErrorType(err.Error()))
				}

				instructionCount = len(p.ProgramStructure.ProgramInstructions)

			case "intsub":

				interruptType := strings.Split(statementValue, " ")[0]
				interruptLabel := strings.Split(statementValue, " ")[1]

				if _, exists := constants.InterruptInts[interruptType]; !exists || constants.InterruptInts[interruptType] < constants.IntMouseMove {
					p.parsingError(ErrSymbol, ErrorType(fmt.Sprintf("unknown interrupt '%s'", interruptType)))
				}

				if !slices.Contains(p.ProgramStructure.LabelNames, interruptLabel) {
					p.parsingError(ErrSymbol, ErrorType(fmt.Sprintf("unrecognized label '%s'", interruptLabel)))
				}

				p.ProgramStructure.InterruptSubscriptions[interruptType] = InterruptSubscription{
					InterruptName: interruptType,
					Interrupt:     constants.InterruptInts[interruptType],
					LabelName:     interruptLabel,
				}

			}

		} else {

			lineSplit := strings.Split(line, " ")

			// Instructions

			if _, exists := constants.InstructionInts[lineSplit[0]]; !exists {
				p.parsingError(ErrDoesNotExist, InstructionDoesNotExist)
			}

			var argCount []int = constants.InstructionArgumentCounts[constants.Instruction(constants.InstructionInts[lineSplit[0]])]

			if !slices.Contains(argCount, len(lineSplit)-1) {

				if len(argCount) > 1 {

					var argList string

					for _, v := range argCount {
						argList += strconv.Itoa(v)
						argList += " "
					}

					p.parsingError(
						ErrWrongNumArgs,
						ErrorType(
							fmt.Sprintf("instruction '%s' expects %s arguments got %d", lineSplit[0], argList, len(lineSplit)-1),
						),
					)

				} else {
					p.parsingError(
						ErrWrongNumArgs,
						ErrorType(
							fmt.Sprintf("too many arguments in call to '%s' - was expecting %d got %d", lineSplit[0], argCount[0], len(lineSplit)-1),
						),
					)
				}

			}

			for _, arg := range lineSplit[1:] {

				// Interrupts

				if lineSplit[0] == "int" {
					if _, exists := constants.InterruptInts[arg]; !exists {
						p.parsingError(ErrSymbol, InvalidArgument)
					}

				} else if arg[0] == '@' && (lineSplit[0] == "lda" || lineSplit[0] == "sta") {

					// Checks definitions for valid argument.

					if !slices.Contains(p.ProgramStructure.DefinitionNames, arg[1:]) {
						p.parsingError(
							ErrDoesNotExist,
							ErrorType(fmt.Sprintf("definition '%s' does not exist", lineSplit[1][1:])),
						)
					}

				} else if lineSplit[0] == "jmp" || lineSplit[0] == "cndjmp" || lineSplit[0] == "call" || lineSplit[0] == "cndcall" {

					if arg[0] == '@' {

						if !slices.Contains(p.ProgramStructure.LabelNames, arg[1:]) {
							p.parsingError(ErrSymbol, SymbolDoesNotExist)
						}

					}

				} else {

					if _, exists := constants.RegisterInts[arg]; !exists {

						p.parsingError(ErrDoesNotExist, ErrorType(fmt.Sprintf("unknown register '%s'", arg)))

					}

				}

			}

			//If does exist, continue

			instructionToBeAdded := Instruction{
				ArgumentCount: uint32(len(strings.Split(line, " ")) - 1),
				StringData:    strings.Split(line, " ")[1:],
				Instruction:   constants.InstructionInts[strings.Split(line, " ")[0]],
			}

			p.ProgramStructure.ProgramInstructions = append(p.ProgramStructure.ProgramInstructions, instructionToBeAdded)

			instructionCount++

		}

	}

	if _, exists := p.ProgramStructure.ProgramLabels["start"]; !exists && !p.Imported {
		p.parsingError(ErrSyntax, "no entrypoint found")
	}

	return p.ProgramStructure, nil
}

// Method for combining the parsers program structure with another program structure.
//
// Used for imports.
func (p *Parser) combine(p1 ProgramStructure) (ProgramStructure, error) {

	var combined ProgramStructure

	//Check for circular imports

	//Combine imports

	combined.ImportedFiles = append(p.ProgramStructure.ImportedFiles[:], p1.ImportedFiles...)

	if slices.Contains(p1.ImportedFiles, p.FileName) {

		p.parsingError(ErrImport, CircularImport)

	}

	// Combine interrupt subscriptions

	combined.InterruptSubscriptions = util.CombineMap[map[string]InterruptSubscription](p.ProgramStructure.InterruptSubscriptions, p1.InterruptSubscriptions)

	//Merge splices & check for name conflicts

	for _, v := range p.ProgramStructure.AllNames {

		if slices.Contains(p1.AllNames, v) {
			return ProgramStructure{}, ErrSymbol
		}

	}

	if len(p1.AllNames) > 0 {
		combined.AllNames = append(p.ProgramStructure.AllNames, p1.AllNames...)
	} else {
		combined.AllNames = p.ProgramStructure.AllNames
	}

	// Update label offsets

	for k, v := range p1.ProgramLabels {
		p1.ProgramLabels[k] = ProgramLabel{
			Name:              v.Name,
			InstructionOffset: v.InstructionOffset + len(p.ProgramStructure.ProgramInstructions),
		}
	}

	if len(p1.ProgramLabels) > 0 {
		combined.ProgramLabels = util.CombineMap[map[string]ProgramLabel](p.ProgramStructure.ProgramLabels, p1.ProgramLabels)
	} else {
		combined.ProgramLabels = p.ProgramStructure.ProgramLabels
	}

	if len(p1.Definitions) > 0 {
		combined.Definitions = util.CombineMap[map[string]Definition](p.ProgramStructure.Definitions, p1.Definitions)
	} else {
		combined.Definitions = p.ProgramStructure.Definitions
	}

	// Combine instructions

	combined.ProgramInstructions = append(p.ProgramStructure.ProgramInstructions, p1.ProgramInstructions...)

	return combined, nil

}

// Name collision function
//
// Checks for any name collisions in the parser, and returns error string.
func (p *Parser) nameCollision(s string) (errMessage string, isCollision bool) {

	isCollision = false

	if _, exists := constants.InstructionInts[s]; exists {
		errMessage = fmt.Sprintf("name %s shares name with instruction", s)
		isCollision = true
	} else if _, exists := constants.RegisterInts[s]; exists {
		errMessage = fmt.Sprintf("name %s shares name with register", s)
		isCollision = true
	} else if _, exists := constants.InterruptInts[s]; exists {
		errMessage = fmt.Sprintf("name %s shares name with interrupt", s)
		isCollision = true
	} else if slices.Contains(p.ProgramStructure.AllNames, s) {
		errMessage = fmt.Sprintf("%s collides with %s", s, s)
		isCollision = true
	}

	return errMessage, isCollision
}
