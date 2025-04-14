package compiler

import (
	"encoding/binary"
	"fmt"
	"log"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"
)

type DisassembledProgram struct {
	ProgramDefinitions [][]byte                       `json:"program_definitions"`
	InterruptTable     map[constants.Interrupt]uint32 `json:"interrupt_table"`
	Instructions       []Instruction                  `json:"instructions"`
	StartIndexes       []uint32                       `json:"start_indexes"`
}

var regMap map[constants.Register]string
var interruptMap map[constants.Interrupt]string

func init() {
	regMap = make(map[constants.Register]string)
	interruptMap = make(map[constants.Interrupt]string)

	for k, v := range constants.RegisterInts {
		regMap[constants.Register(v)] = k
	}

	for k, v := range constants.InterruptInts {
		interruptMap[v] = k
	}
}

// Decodes individual instructions.
func decodeInstruction(b []byte) Instruction {

	i := Instruction{}

	itn := constants.Instruction(b[0])
	itnDataBytes := b[1:]

	var itnData []uint32

	// Get instruction arguments
	for i := 0; i < constants.InstructionArgumentCounts[itn]; i += 2 {
		itnData = append(itnData, uint32(binary.LittleEndian.Uint16(itnDataBytes[i:i+2])))
	}

	argumentData := ""

	for _, v := range itnData {

		// If the instruction is one where a memory address is passed as an arguement instead of a register.
		if itn == constants.ILoad ||
			itn == constants.IStore ||
			itn == constants.IJump ||
			itn == constants.IConditionalJump ||
			itn == constants.ICall ||
			itn == constants.IConditionalCall {

			argumentData = fmt.Sprintf("0x"+"%08X", itnData[0])

		} else {

			if itn == constants.ICallInterrupt {
				argumentData = interruptMap[constants.Interrupt(v)]
			} else {
				argumentData = regMap[constants.Register(v)]
			}
		}

		i.StringData = append(i.StringData, argumentData)
	}

	i.Instruction = uint32(itn)

	return i

}

// Main disassemble method.
//
// Takes program bytes and returns a DisassembledProgram struct.
// The disassembled struct is similar to an assembled program struct but it is missing many of the fields that the assembled program struct has.
func Disassemble(programBytes []byte, verbose bool) (DisassembledProgram, error) {

	var dataBlockStart uint32
	var dataBlockBytes []byte

	var interruptTableStart uint32
	var interruptBlockBytes []byte

	var instructionStart uint32
	var instructionEntryPoint uint32
	var instructionBytes []byte

	var program DisassembledProgram

	//Extract header and

	if verbose {
		log.Printf("Got program with %d byte(s)", len(programBytes)-4)
	}

	programBytes = programBytes[4:]

	program.StartIndexes = make([]uint32, 4)

	interruptTableStart = binary.LittleEndian.Uint32(programBytes[:4])
	program.StartIndexes[0] = interruptTableStart
	dataBlockStart = binary.LittleEndian.Uint32(programBytes[4:8])
	program.StartIndexes[1] = dataBlockStart
	instructionStart = binary.LittleEndian.Uint32(programBytes[8:12])
	program.StartIndexes[2] = instructionStart
	instructionEntryPoint = binary.LittleEndian.Uint32(programBytes[12:16])
	program.StartIndexes[3] = instructionEntryPoint

	if verbose {

		log.Println("Indexes are:")
		log.Printf("Data block start index %d", dataBlockStart)
		log.Printf("Interrupt table start index %d", interruptTableStart)
		log.Printf("Instruction start index %d", instructionStart)
		log.Printf("Instruction entry point %d", instructionEntryPoint)

	}

	dataBlockBytes = programBytes[dataBlockStart:instructionStart]
	interruptBlockBytes = programBytes[interruptTableStart:dataBlockStart]
	instructionBytes = programBytes[instructionStart : len(programBytes)-int(PadSize)]

	if verbose {
		log.Println("Disassembling data table...")
	}

	//----------------
	//Build data table
	//----------------

	program.ProgramDefinitions = make([][]byte, 0)

	fmt.Println(dataBlockBytes)

	//Break definitions up into individual byte arrays

	var i uint32 = 0
	var definitionLength uint32 = 0
	for i = 0; i < uint32(len(dataBlockBytes)); i += definitionLength + 4 {

		definitionLength = binary.LittleEndian.Uint32(dataBlockBytes[i : i+4])

		if verbose {
			log.Printf("value of i = %d", i)
			log.Printf("%d definition - %d byte(s) long.", len(program.ProgramDefinitions)-1, definitionLength)
		}

		program.ProgramDefinitions = append(
			program.ProgramDefinitions,
			dataBlockBytes[i+4:i+4+definitionLength],
		)

	}

	if verbose {
		log.Printf("Found %d definition(s)", len(program.ProgramDefinitions))
		log.Println("Disassembling interrupt table...")
	}

	// Build interrupt table

	program.InterruptTable = make(map[constants.Interrupt]uint32)

	if verbose {
		log.Printf("Interrupt byte length %d", len(interruptBlockBytes))
	}

	for _, v := range util.SliceChunks(interruptBlockBytes, 6) {

		interrupt := constants.Interrupt(binary.LittleEndian.Uint16(v[:2]))
		labelAddr := binary.LittleEndian.Uint32(v[2:])

		program.InterruptTable[interrupt] = labelAddr

	}

	if verbose {
		log.Println("Disassembling instructions...")
	}

	// Decode instructions

	for _, v := range util.SliceChunks(instructionBytes, int(InstructionLength)) {

		program.Instructions = append(program.Instructions, decodeInstruction(v))

	}

	return program, nil

}

func allZero(b []byte) bool {

	for _, v := range b {

		if v != 0 {
			return false
		}

	}

	return true

}
