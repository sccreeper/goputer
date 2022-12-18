package compiler

//Handles catches errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type ErrorType string

const (
	UnexpectedEndStatement  ErrorType = "unexpected end statement"
	NestingError            ErrorType = "cannot nest jump blocks"
	MinimumNameLength       ErrorType = "minimum length of one"
	NameConflict            ErrorType = "symbol already exists"
	InstructionDoesNotExist ErrorType = "instruction does not exist"
	SymbolDoesNotExist      ErrorType = "symbol does not exist"
	InvalidArgument         ErrorType = "invalid argument"

	CircularImport ErrorType = "circular import"
	FileNotFound   ErrorType = "file not found"
)

var ErrSyntax error = errors.New("syntax error")
var ErrSymbol error = errors.New("symbol error")
var ErrDoesNotExist error = errors.New("does not exist error")
var ErrInvalidArgument error = errors.New("invalid argument")
var ErrImport error = errors.New("import error")

var RedError color.Color = *color.New(color.FgHiRed, color.Bold)
var ItalicCode color.Color = *color.New(color.Italic)

// Handles a parsing error
func (p *Parser) parsing_error(e error, error_type ErrorType) {

	var error_text string

	error_text += "Error\n"

	if p.LineIndex != -1 {
		error_text += fmt.Sprintf("In file '%s' on line %d\n", p.FileName, p.LineIndex+1)
	} else {
		error_text += fmt.Sprintf("In file '%s'\n", p.FileName)
	}

	error_text += format_line(p.CurrentLine) + "\n"

	switch e {
	case ErrSyntax:
		error_text += fmt.Sprintf("%s %s\n", RedError.Sprint("Syntax error:"), error_type)
	case ErrSymbol:
		error_text += fmt.Sprintf("%s %s\n", RedError.Sprint("Symbol error:"), error_type)
	case ErrDoesNotExist:
		error_text += fmt.Sprintf("%s %s\n", RedError.Sprint("Does not exist:"), error_type)
	case ErrInvalidArgument:
		error_text += fmt.Sprintf("%s %s\n", RedError.Sprint("Invalid argument:"), error_type)
	case ErrImport:
		error_text += fmt.Sprintf("%s %s\n", RedError.Sprint("Import error:"), error_type)
	}

	p.ErrorHandler(error_type, error_text)

}

func format_line(line string) string {

	line_data := strings.Split(line, " ")

	if line_data[0] == "def" {

		def_line := []string{}
		def_line = append(def_line, "def", line_data[1])

		if line_data[2][0] == '"' {

			def_line = append(def_line, line[len("def")+len(line_data[1])-1+3:])

		} else {
			def_line = append(def_line, line_data[2])
		}

		line_data = def_line

	}

	var final_line string = ""

	for index, v := range line_data {

		if index == 0 {

			final_line += ItalicCode.Sprint(color.GreenString(v))

		} else {

			final_line += ItalicCode.Sprint(color.CyanString(v))

		}

		final_line += " "

	}

	return final_line

}
