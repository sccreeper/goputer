package compiler

import (
	"encoding/binary"
	"fmt"
	"log"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"

	"golang.org/x/exp/slices"
)

type DisassembledProgram struct {
	ProgramDefinitions [][]byte                       `json:"program_definitions"`
	InterruptTable     map[constants.Interrupt]uint32 `json:"interrupt_table"`
	JumpBlocks         map[uint32][]Instruction       `json:"jump_blocks"`
	Instructions       []Instruction                  `json:"instructions"`
	StartIndexes       []uint32                       `json:"start_indexes"`
}

// Decodes individual instructions.
func decode_instruction(b []byte) Instruction {

	i := Instruction{}

	itn := constants.Instruction(b[0])
	itn_data_bytes := b[1:]

	var single_data bool
	var itn_data []uint32

	//Reverse maps

	reg_map := make(map[constants.Register]string)

	for k, v := range constants.RegisterInts {

		reg_map[constants.Register(v)] = k

	}

	interrupt_map := make(map[constants.Interrupt]string)

	for k, v := range constants.InterruptInts {

		interrupt_map[v] = k

	}

	//Determine if instruction is two args or single arg
	if slices.Contains(constants.SingleArgInstructions, itn) {

		itn_data = append(itn_data, binary.LittleEndian.Uint32(itn_data_bytes))

		single_data = true
	} else {
		itn_data = append(itn_data, uint32(binary.LittleEndian.Uint16(itn_data_bytes[:2])))
		itn_data = append(itn_data, uint32(binary.LittleEndian.Uint16(itn_data_bytes[2:4])))

		single_data = false
	}

	d := ""

	for _, v := range itn_data {

		if itn == constants.ILoad || itn == constants.IStore || itn == constants.IJump || itn == constants.IConditionalJump {

			d = fmt.Sprintf("0x"+"%08X", itn_data[0])

		} else {

			if itn == constants.ICallInterrupt {
				d = interrupt_map[constants.Interrupt(v)]
			} else {
				d = reg_map[constants.Register(v)]
			}
		}

		i.Data = append(i.Data, d)
	}

	i.SingleData = single_data
	i.Instruction = uint32(itn)

	return i

}

// Main disassemble method.
//
// Takes program bytes and returns a DisassembledProgram struct.
// The disassembled struct is similar to an assembled program struct but it is missing many of the fields that the assembled program struct has.
func Disassemble(program_bytes []byte, verbose bool) (DisassembledProgram, error) {

	var data_block_start uint32
	var data_block_bytes []byte

	var interrupt_table_start uint32
	var interrupt_block_bytes []byte

	var jump_block_start uint32
	var jump_block_bytes []byte

	var instruction_start uint32
	var instruction_bytes []byte

	var program DisassembledProgram

	var byte_index uint32 = 0

	//Extract header and

	if verbose {
		log.Printf("Got program with %d byte(s)", len(program_bytes)-4)
	}

	program_bytes = program_bytes[4:]

	program.StartIndexes = make([]uint32, 4)

	data_block_start = binary.LittleEndian.Uint32(program_bytes[:4])
	program.StartIndexes[0] = data_block_start
	jump_block_start = binary.LittleEndian.Uint32(program_bytes[4:8])
	program.StartIndexes[1] = jump_block_start
	interrupt_table_start = binary.LittleEndian.Uint32(program_bytes[8:12])
	program.StartIndexes[2] = interrupt_table_start
	instruction_start = binary.LittleEndian.Uint32(program_bytes[12:16])
	program.StartIndexes[3] = instruction_start

	if verbose {

		log.Println("Indexes are:")
		log.Printf("Data block start index %d", data_block_start)
		log.Printf("Jump block start index %d", jump_block_start)
		log.Printf("Interrupt table start index %d", interrupt_table_start)
		log.Printf("Instruction start index %d", instruction_start)

	}

	data_block_bytes = program_bytes[data_block_start : jump_block_start-PadSize]
	jump_block_bytes = program_bytes[jump_block_start : interrupt_table_start-PadSize]
	interrupt_block_bytes = program_bytes[interrupt_table_start : instruction_start-PadSize]
	instruction_bytes = program_bytes[instruction_start : len(program_bytes)-int(PadSize)]

	if verbose {
		log.Println("Disassembling data table...")
	}

	//----------------
	//Build data table
	//----------------

	byte_index = data_block_start

	program.ProgramDefinitions = append(program.ProgramDefinitions, []byte{})
	var definition_index uint32 = 0
	var definition_length_index uint32 = 0
	var definition_length uint32 = 0
	var in_definition bool = false
	var data_bytes_index uint32 = 0

	//Break definitions up into individual byte arrays

	for range data_block_bytes {

		if int(data_bytes_index) > len(data_block_bytes) || int(data_bytes_index+4) > len(data_block_bytes) {
			break
		}

		if !in_definition {

			definition_length = binary.LittleEndian.Uint32(data_block_bytes[data_bytes_index : data_bytes_index+4])

			if verbose {
				log.Printf("%d definition - %d byte(s) long.", definition_index, definition_length)
			}

			data_bytes_index += 4
			in_definition = true
			continue

		} else if definition_length_index >= definition_length {

			in_definition = false

			program.ProgramDefinitions = append(program.ProgramDefinitions, []byte{})

			definition_index++
			definition_length_index = 0

			continue
		} else {

			program.ProgramDefinitions[definition_index] = append(program.ProgramDefinitions[definition_index], data_block_bytes[data_bytes_index])
			definition_length_index++

		}

		data_bytes_index++

	}

	if verbose {
		log.Println("Disassembling jump blocks...")
	}

	//---------------------
	//Build jump blocks
	//--------------------

	byte_index = jump_block_start

	program.JumpBlocks = make(map[uint32][]Instruction)

	jump_block_addr_index := byte_index

	in_jump_block := false

	for _, v := range util.SliceChunks(jump_block_bytes, int(InstructionLength)) {

		if !in_jump_block {
			program.JumpBlocks[jump_block_addr_index] = make([]Instruction, 0)

			program.JumpBlocks[jump_block_addr_index] = append(program.JumpBlocks[jump_block_addr_index], decode_instruction(v))

			in_jump_block = true
		} else if all_zero(v) {
			in_jump_block = false

			jump_block_addr_index += uint32((len(program.JumpBlocks[jump_block_addr_index]) * int(InstructionLength)) + int(InstructionLength))
		} else {
			program.JumpBlocks[jump_block_addr_index] = append(program.JumpBlocks[jump_block_addr_index], decode_instruction(v))
		}

	}

	if verbose {
		log.Printf("Found %d jump block(s)", len(program.JumpBlocks))
		log.Println("Disassembling interrupt table...")
	}

	//---------------------
	//Build interrupt table
	//---------------------

	byte_index = interrupt_table_start

	//Reverse interrupt map

	interrupt_map := make(map[constants.Interrupt]string)

	for k, v := range constants.SubscribableInterrupts {

		interrupt_map[v] = k

	}

	//Decode bytes

	program.InterruptTable = make(map[constants.Interrupt]uint32)

	if verbose {
		log.Printf("Interrupt byte length %d", len(interrupt_block_bytes))
	}

	for _, v := range util.SliceChunks(interrupt_block_bytes, 6) {

		//log.Println(current_bytes)

		interrupt := constants.Interrupt(binary.LittleEndian.Uint16(v[:2]))
		jump_block_addr := binary.LittleEndian.Uint32(v[2:])

		program.InterruptTable[interrupt] = jump_block_addr

	}

	if verbose {
		log.Println("Disassembling instructions...")
	}

	//---------------------------
	//Decode instructions
	//---------------------------

	for _, v := range util.SliceChunks(instruction_bytes, int(InstructionLength)) {

		program.Instructions = append(program.Instructions, decode_instruction(v))

	}

	return program, nil

}

func all_zero(b []byte) bool {

	for _, v := range b {

		if v != 0 {
			return false
		}

	}

	return true

}
