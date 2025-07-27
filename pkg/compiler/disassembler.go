package compiler

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	c "sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"
	"strings"
)

type DisassembledProgram struct {
	ProgramDefinitions [][]byte                       `json:"program_definitions"`
	InterruptTable     map[c.Interrupt]uint32 `json:"interrupt_table"`
	Instructions       []Instruction                  `json:"instructions"`
	StartIndexes       []uint32                       `json:"start_indexes"`
}

var regMap map[c.Register]string
var interruptMap map[c.Interrupt]string
var instructionMap map[c.Instruction]string

func init() {
	regMap = make(map[c.Register]string)
	interruptMap = make(map[c.Interrupt]string)
	instructionMap = make(map[c.Instruction]string)

	for k, v := range c.RegisterInts {
		regMap[c.Register(v)] = k
	}

	for k, v := range c.InterruptInts {
		interruptMap[v] = k
	}

	for k, v := range c.InstructionInts {
		instructionMap[c.Instruction(v)] = k
	}
}

func DecodeInstructionString(b []byte) (string, error) {
	
	itn, err := DecodeInstruction(b)

	return fmt.Sprintf("%s %s", instructionMap[c.Instruction(itn.Instruction)], strings.Join(itn.StringData, " ")), err

}

// Decodes individual instructions.
func DecodeInstruction(b []byte) (Instruction, error) {

	i := Instruction{}

	itn := c.Instruction(b[0]) & c.Instruction(c.InstructionMask)
	itnDataBytes := b[1:]

	var itnData []uint32

	if int(itn) > len(c.InstructionArgumentCounts) {
		fmt.Println(itn)
		return Instruction{}, errors.New("error decoding instruction")
	}

	// Get instruction arguments
	for i := range c.InstructionArgumentCounts[itn][0] {
		itnData = append(itnData, uint32(binary.LittleEndian.Uint16(itnDataBytes[(i*2):(i*2)+2])))
	}

	argumentData := ""

	if b[0] & byte(c.ItnFlagLongArgImmediate) == byte(c.ItnFlagLongArgImmediate) {
		
		argumentData = fmt.Sprintf("$%d", binary.LittleEndian.Uint32(itnDataBytes))

		i.StringData = append(i.StringData, argumentData)

	} else if (b[0] & byte(c.ItnFlagLeftArgImmediate)) != 0 || (b[0] & byte(c.ItnFlagRightArgImmediate)) != 0 {
		
		immediateValue := binary.LittleEndian.Uint32(itnDataBytes) & c.InstructionArgImmediateMask
		immediateRegister := (binary.LittleEndian.Uint32(itnDataBytes) & c.InstructionArgRegisterMask) >> 26

		if (b[0] & byte(c.ItnFlagLeftArgImmediate)) != 0 {
			argumentData = fmt.Sprintf("$%d %s", immediateValue, regMap[c.Register(immediateRegister)])
		} else {
			argumentData = fmt.Sprintf("%s $%d", regMap[c.Register(immediateRegister)], immediateValue)
		}

		i.StringData = append(i.StringData, argumentData)

	} else {
		for _, v := range itnData {

			// If the instruction is one where a memory address is passed as an argument instead of a register.
			if itn == c.ILoad ||
				itn == c.IStore ||
				itn == c.IJump ||
				itn == c.IConditionalJump ||
				itn == c.ICall ||
				itn == c.IConditionalCall {

				argumentData = fmt.Sprintf("0x%08X", binary.LittleEndian.Uint32(itnDataBytes))

			} else {

				if itn == c.ICallInterrupt {
					argumentData = interruptMap[c.Interrupt(v)]
				} else {
					argumentData = regMap[c.Register(v)]
				}
			}

			i.StringData = append(i.StringData, argumentData)
		}
	}

	i.Instruction = uint32(itn)

	return i, nil

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

	program.InterruptTable = make(map[c.Interrupt]uint32)

	if verbose {
		log.Printf("Interrupt byte length %d", len(interruptBlockBytes))
	}

	for _, v := range util.SliceChunks(interruptBlockBytes, 6) {

		interrupt := c.Interrupt(binary.LittleEndian.Uint16(v[:2]))
		labelAddr := binary.LittleEndian.Uint32(v[2:])

		program.InterruptTable[interrupt] = labelAddr

	}

	if verbose {
		log.Println("Disassembling instructions...")
	}

	// Decode instructions

	for _, v := range util.SliceChunks(instructionBytes, int(InstructionLength)) {

		itn, err := DecodeInstruction(v)
		if err != nil {
			return DisassembledProgram{}, err
		}

		program.Instructions = append(program.Instructions, itn)

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
