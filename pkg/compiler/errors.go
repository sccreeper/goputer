package compiler

//Handles catches errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type ErrorMessage string

const (
	UnexpectedEndStatement  ErrorMessage = "unexpected end statement" //Error for a end statement outside of a jump block.
	NestingError            ErrorMessage = "cannot nest jump blocks"  //Error for nesting jump blocks.
	MinimumNameLength       ErrorMessage = "minimum length of one"    //Error for having a def/jump block name which isn't at least 1 character long.
	NameConflict            ErrorMessage = "symbol already exists"
	InstructionDoesNotExist ErrorMessage = "instruction does not exist"
	SymbolDoesNotExist      ErrorMessage = "symbol does not exist" //When a symbol for a definition or a jump block doesn't exist.
	InvalidArgument         ErrorMessage = "invalid argument"
	InvalidValue            ErrorMessage = "invalid value"

	CircularImport ErrorMessage = "circular import"
	FileNotFound   ErrorMessage = "file not found"
)

var ErrSyntax error = errors.New("syntax error")
var ErrValue error = errors.New("value error")
var ErrSymbol error = errors.New("symbol error")
var ErrDoesNotExist error = errors.New("does not exist error")
var ErrInvalidArgument error = errors.New("invalid argument")
var ErrImport error = errors.New("import error")
var ErrWrongNumArgs error = errors.New("wrong number of arguments")
var ErrFile = errors.New("error whilst reading file")

var RedError color.Color = *color.New(color.FgHiRed, color.Bold)
var ItalicCode color.Color = *color.New(color.Italic)

// Handles a parsing error
func (p *Parser) parsingError(e error, errorType ErrorMessage) {

	var errorText string

	errorText += "Error\n"

	if p.LineIndex != -1 {
		errorText += fmt.Sprintf("In file '%s' on line %d\n", p.FileName, p.LineIndex+1)
	} else {
		errorText += fmt.Sprintf("In file '%s'\n", p.FileName)
	}

	errorText += formatLine(p.CurrentLine) + "\n"

	switch e {
	case ErrSyntax:
		errorText += fmt.Sprintf("%s %s\n", RedError.Sprint("Syntax error:"), errorType)
	case ErrSymbol:
		errorText += fmt.Sprintf("%s %s\n", RedError.Sprint("Symbol error:"), errorType)
	case ErrDoesNotExist:
		errorText += fmt.Sprintf("%s %s\n", RedError.Sprint("Does not exist:"), errorType)
	case ErrInvalidArgument:
		errorText += fmt.Sprintf("%s %s\n", RedError.Sprint("Invalid argument:"), errorType)
	case ErrImport:
		errorText += fmt.Sprintf("%s %s\n", RedError.Sprint("Import error:"), errorType)
	case ErrWrongNumArgs:
		errorText += fmt.Sprintf("%s %s\n", RedError.Sprint("Wrong number of arguments:"), errorType)
	default:
		errorText += string(errorType)
	}

	p.ErrorHandler(errorType, errorText)

}

func formatLine(line string) string {

	lineData := strings.Split(line, " ")

	if lineData[0] == "def" {

		defLine := []string{}
		defLine = append(defLine, "def", lineData[1])

		if lineData[2][0] == '"' {

			defLine = append(defLine, line[len("def")+len(lineData[1])-1+3:])

		} else {
			defLine = append(defLine, lineData[2])
		}

		lineData = defLine

	}

	var finalLine string = ""

	for index, v := range lineData {

		if index == 0 {

			finalLine += ItalicCode.Sprint(color.GreenString(v))

		} else {

			finalLine += ItalicCode.Sprint(color.CyanString(v))

		}

		finalLine += " "

	}

	return finalLine

}
