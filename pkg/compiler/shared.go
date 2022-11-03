package compiler

import "sccreeper/govm/pkg/constants"

//File for code shared in the compiler

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
	InstructionLength uint32 = 5    //Instruction length in bytes
	BlockAddrSize     uint32 = 16   // Size of the block address header
	PadSize           uint32 = 4    //Padding size inbetween blocks
	PadValue          byte   = 0xFF //Value to pad blocks with
	InterruptLength   uint32 = 6    //Length of interrupt in bytes (1 uint16, 1 uint32)
	StackSize         uint32 = 1024 //Default stack size (256 uint32)
)
