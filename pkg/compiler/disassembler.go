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
	JumpBlocks         map[uint32][]Instruction       `json:"jump_blocks"`
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

	var jumpBlockStart uint32
	var jumpBlockBytes []byte

	var instructionStart uint32
	var instructionBytes []byte

	var program DisassembledProgram

	var byteIndex uint32 = 0

	//Extract header and

	if verbose {
		log.Printf("Got program with %d byte(s)", len(programBytes)-4)
	}

	programBytes = programBytes[4:]

	program.StartIndexes = make([]uint32, 4)

	dataBlockStart = binary.LittleEndian.Uint32(programBytes[:4])
	program.StartIndexes[0] = dataBlockStart
	jumpBlockStart = binary.LittleEndian.Uint32(programBytes[4:8])
	program.StartIndexes[1] = jumpBlockStart
	interruptTableStart = binary.LittleEndian.Uint32(programBytes[8:12])
	program.StartIndexes[2] = interruptTableStart
	instructionStart = binary.LittleEndian.Uint32(programBytes[12:16])
	program.StartIndexes[3] = instructionStart

	if verbose {

		log.Println("Indexes are:")
		log.Printf("Data block start index %d", dataBlockStart)
		log.Printf("Jump block start index %d", jumpBlockStart)
		log.Printf("Interrupt table start index %d", interruptTableStart)
		log.Printf("Instruction start index %d", instructionStart)

	}

	dataBlockBytes = programBytes[dataBlockStart : jumpBlockStart-PadSize]
	jumpBlockBytes = programBytes[jumpBlockStart : interruptTableStart-PadSize]
	interruptBlockBytes = programBytes[interruptTableStart : instructionStart-PadSize]
	instructionBytes = programBytes[instructionStart : len(programBytes)-int(PadSize)]

	if verbose {
		log.Println("Disassembling data table...")
	}

	//----------------
	//Build data table
	//----------------

	byteIndex = dataBlockStart

	program.ProgramDefinitions = append(program.ProgramDefinitions, []byte{})
	var definitionIndex uint32 = 0
	var definitionLengthIndex uint32 = 0
	var definitionLength uint32 = 0
	var inDefinition bool = false
	var dataBytesIndex uint32 = 0

	//Break definitions up into individual byte arrays

	for range dataBlockBytes {

		if int(dataBytesIndex) > len(dataBlockBytes) || int(dataBytesIndex+4) > len(dataBlockBytes) {
			break
		}

		if !inDefinition {

			definitionLength = binary.LittleEndian.Uint32(dataBlockBytes[dataBytesIndex : dataBytesIndex+4])

			if verbose {
				log.Printf("%d definition - %d byte(s) long.", definitionIndex, definitionLength)
			}

			dataBytesIndex += 4
			inDefinition = true
			continue

		} else if definitionLengthIndex >= definitionLength {

			inDefinition = false

			program.ProgramDefinitions = append(program.ProgramDefinitions, []byte{})

			definitionIndex++
			definitionLengthIndex = 0

			continue
		} else {

			program.ProgramDefinitions[definitionIndex] = append(program.ProgramDefinitions[definitionIndex], dataBlockBytes[dataBytesIndex])
			definitionLengthIndex++

		}

		dataBytesIndex++

	}

	if verbose {
		log.Println("Disassembling jump blocks...")
	}

	//---------------------
	//Build jump blocks
	//--------------------

	byteIndex = jumpBlockStart

	program.JumpBlocks = make(map[uint32][]Instruction)

	jumpBlockAddrIndex := byteIndex

	inJumpBlock := false

	for _, v := range util.SliceChunks(jumpBlockBytes, int(InstructionLength)) {

		if !inJumpBlock {
			program.JumpBlocks[jumpBlockAddrIndex] = make([]Instruction, 0)

			program.JumpBlocks[jumpBlockAddrIndex] = append(program.JumpBlocks[jumpBlockAddrIndex], decodeInstruction(v))

			inJumpBlock = true
		} else if allZero(v) {
			inJumpBlock = false

			jumpBlockAddrIndex += uint32((len(program.JumpBlocks[jumpBlockAddrIndex]) * int(InstructionLength)) + int(InstructionLength))
		} else {
			program.JumpBlocks[jumpBlockAddrIndex] = append(program.JumpBlocks[jumpBlockAddrIndex], decodeInstruction(v))
		}

	}

	if verbose {
		log.Printf("Found %d jump block(s)", len(program.JumpBlocks))
		log.Println("Disassembling interrupt table...")
	}

	// Build interrupt table

	byteIndex = interruptTableStart

	program.InterruptTable = make(map[constants.Interrupt]uint32)

	if verbose {
		log.Printf("Interrupt byte length %d", len(interruptBlockBytes))
	}

	for _, v := range util.SliceChunks(interruptBlockBytes, 6) {

		interrupt := constants.Interrupt(binary.LittleEndian.Uint16(v[:2]))
		jumpBlockAddr := binary.LittleEndian.Uint32(v[2:])

		program.InterruptTable[interrupt] = jumpBlockAddr

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
