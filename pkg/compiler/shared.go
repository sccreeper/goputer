package compiler

import "sccreeper/goputer/pkg/constants"

//File for code shared in the compiler

// Compiler backend

//Config for CLI

type CompilerConfig struct {
	OutputPath string
	FilePath   string

	OutputJSON bool
	JSONPath   string
	Verbose    bool
}

// Types for statements
type Instruction struct {
	ArgumentCount uint32   `json:"argument_count"`
	StringData    []string `json:"data"`
	ByteData      []byte   `json:"byte_data"`

	Instruction uint32 `json:"instruction"`
}

type Definition struct {
	Name       string `json:"name"`
	StringData string `json:"string_data"`
	ByteData   []byte `json:"byte_data"`
	Address    uint32 `json:"address"`

	Type constants.DefType `json:"type"`
}

type InterruptSubscription struct {
	InterruptName string              `json:"interrupt_name"`
	Interrupt     constants.Interrupt `json:"interrupt"`
	LabelName     string              `json:"label_name"`
}

type CodeBlock struct {
	Name         string        `json:"name"`
	Instructions []Instruction `json:"instructions"`
}

// Struct for holding program data after first stage of parsing
type ProgramStructure struct {
	AllNames      []string `json:"all_names"`
	ImportedFiles []string `json:"imported_files"`

	// where string = interrupt type/name
	InterruptSubscriptions map[string]InterruptSubscription `json:"interrupt_subscriptions"`
	// where string = name of program label
	ProgramLabels       map[string]ProgramLabel `json:"program_labels"`
	ProgramInstructions []Instruction           `json:"program_instructions"`
	// where definition = name of program label
	Definitions map[string]Definition `json:"definitions"`
}

type ProgramLabel struct {
	Name              string `json:"name"`
	InstructionOffset int    `json:"instruction_offset"`
}

// File used for compilers filesystem.
type VFSFile struct {
	RealPath   string
	FakePath   string
	Name       string
	Data       []byte
	StringData string
}

//Constants

const (
	InstructionLength uint32 = 5       //Instruction length in bytes
	HeaderSize        uint32 = 16      // Size of the header (without magic string)
	PadSize           uint32 = 4       //Padding size inbetween blocks
	PadValue          byte   = 0xFF    //Value to pad blocks with
	InterruptLength   uint32 = 6       //Length of interrupt in bytes (1 uint16, 1 uint32)
	DataStackSize     uint32 = 4 * 256 //Default stack size, 1024 bytes(256 uint32)
	CallStackSize     uint32 = 4 * 128 //Default call stack size, 256 bytes, 64 uint32
	StackSize         uint32 = DataStackSize + CallStackSize
)
