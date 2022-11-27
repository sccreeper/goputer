package compiler

import "sccreeper/goputer/pkg/constants"

//File for code shared in the compiler

// Compiler backend

//Config for CLI

type CompilerConfig struct {
	OutputPath string
	FileName   string

	OutputJSON bool
	JSONPath   string
	Verbose    bool
}

// Types for statements
type Instruction struct {
	SingleData bool     `json:"single_data"`
	Data       []string `json:"data"`

	Instruction uint32 `json:"instruction"`
}

type Definition struct {
	Name string            `json:"name"`
	Data []byte            `json:"data"`
	Type constants.DefType `json:"type"`
}

type InterruptSubscription struct {
	InterruptName string              `json:"interrupt_name"`
	Interrupt     constants.Interrupt `json:"interrupt"`
	JumpBlockName string              `json:"jump_block_name"`
}

type CodeBlock struct {
	Name         string        `json:"name"`
	Instructions []Instruction `json:"instructions"`
}

// Struct for holding program data
type ProgramStructure struct {
	AllNames []string `json:"all_names"`

	InstructionBlockNames  []string                `json:"instruction_block_names"`
	DefNames               []string                `json:"definition_names"`
	InterruptSubscriptions []InterruptSubscription `json:"interrupt_subscriptions"`

	ProgramInstructions []Instruction        `json:"program_instructions"`
	Definitions         []Definition         `json:"definitions"`
	InstructionBlocks   map[string]CodeBlock `json:"instruction_blocks"`
}

//Constants

const (
	InstructionLength uint32 = 5       //Instruction length in bytes
	BlockAddrSize     uint32 = 16      // Size of the block address header
	PadSize           uint32 = 4       //Padding size inbetween blocks
	PadValue          byte   = 0xFF    //Value to pad blocks with
	InterruptLength   uint32 = 6       //Length of interrupt in bytes (1 uint16, 1 uint32)
	DataStackSize     uint32 = 4 * 256 //Default stack size, 1024 bytes(256 uint32)
	CallStackSize     uint32 = 4 * 64  //Default call stack size, 256 bytes, 64 uint32
	StackSize         uint32 = DataStackSize + CallStackSize
)
