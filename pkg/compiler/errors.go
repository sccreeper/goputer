package compiler

//Handles catches errors

import (
	"errors"
	"fmt"
	"os"
)

type ErrorType string

const (
	UnexpectedEndStatement  ErrorType = "unexpected end statement"
	NestingError            ErrorType = "cannot nest jump blocks"
	MinimumNameLength       ErrorType = "minimum length error"
	NameConflict            ErrorType = "symbol already exists"
	InstructionDoesNotExist ErrorType = "instruction does not exist"
	SymbolDoesNotExist      ErrorType = "symbol does not exist"
)

var ErrSyntax error = errors.New("syntax error")
var ErrSymbol error = errors.New("symbol error")
var ErrDoesNotExist error = errors.New("does not exist error")

// Handles a parsing error
func parsing_error(e error, line_number int, file string, line string, error_type ErrorType) {

	if line_number != -1 {
		fmt.Printf("Error in file '%s' on line %d\n", file, line_number+1)
	} else {
		fmt.Printf("Error in file '%s'\n", file)
	}

	fmt.Println(line)

	switch e {
	case ErrSyntax:
		fmt.Printf("Syntax error: %s\n", error_type)
		os.Exit(1)
	case ErrSymbol:
		fmt.Printf("Symbol error: %s\n", error_type)
		os.Exit(1)
	case ErrDoesNotExist:
		fmt.Printf("Does not exist: %s\n", error_type)
		os.Exit(1)
	}

}
