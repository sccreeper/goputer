package compiler

import (
	"encoding/binary"
	"sccreeper/govm/pkg/constants"
	"sccreeper/govm/pkg/util"
)

func generate_bytecode(p ProgramStructure) []byte {

	byte_index := BlockAddrSize + 4
	final_bytes := make([]byte, 0)

	jump_block_addr := make(map[string]uint32)
	data_block_addr := make(map[string]uint32)

	//Generate definition bytecode first

	definition_bytes := []byte{}
	definition_start_index := byte_index
	definition_addr_index := definition_start_index

	for _, d := range p.Definitions {

		data_block_addr[d.Name] = definition_addr_index + StackSize

		length_bytes := make([]byte, 4)

		binary.LittleEndian.PutUint32(length_bytes, uint32(len(d.Data)))

		data_bytes := d.Data

		definition_bytes = append(definition_bytes, length_bytes...)
		definition_bytes = append(definition_bytes, data_bytes...)

		definition_addr_index += uint32(len(length_bytes) + len(data_bytes))

	}

	//Increment the byte index

	byte_index += uint32(len(definition_bytes)) + PadSize

	//Generate the jump block block

	jump_block_start_index := byte_index
	jump_block_addr_index := jump_block_start_index

	jump_block_bytes := []byte{}

	//Order keys first

	for _, s := range p.InstructionBlockNames {

		v := p.InstructionBlocks[s]

		current_jump_block_bytes := []byte{}

		jump_block_addr[v.Name] = jump_block_addr_index

		for _, i := range v.Instructions {

			current_jump_block_bytes = append(current_jump_block_bytes, generate_instruction_bytecode(i, data_block_addr, jump_block_addr)...)

		}

		current_jump_block_bytes = append(current_jump_block_bytes, []byte{0, 0, 0, 0, 0}...)

		jump_block_bytes = append(jump_block_bytes, current_jump_block_bytes...)

		jump_block_addr_index += uint32(len(current_jump_block_bytes))

		current_jump_block_bytes = nil

	}

	byte_index += uint32(len(jump_block_bytes)) + PadSize

	//Generate the interrupt table

	included_interrupts := []constants.Interrupt{}
	current_interrupt_bytes := []byte{}
	interrupt_bytes := []byte{}

	interrupt_block_start_index := byte_index

	for _, v := range p.InterruptSubscriptions {

		i := v.Interrupt
		j_addr := jump_block_addr[v.JumpBlockName]

		i_bytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(i_bytes[:], uint16(i))
		j_addr_bytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(j_addr_bytes[:], j_addr)

		current_interrupt_bytes = append(current_interrupt_bytes, i_bytes...)
		current_interrupt_bytes = append(current_interrupt_bytes, j_addr_bytes...)

		interrupt_bytes = append(interrupt_bytes, current_interrupt_bytes...)

		current_interrupt_bytes = nil

		included_interrupts = append(included_interrupts, v.Interrupt)

	}

	//Generate the rest of the interrupt table bytes

	for _, v := range constants.SubscribableInterrupts {

		if util.SliceContains(included_interrupts, v) {
			continue
		} else {
			i := v
			j_addr := uint32(0x00000000)

			i_bytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(i_bytes[:], uint16(i))
			j_addr_bytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(j_addr_bytes[:], j_addr)

			current_interrupt_bytes = append(current_interrupt_bytes, i_bytes...)
			current_interrupt_bytes = append(current_interrupt_bytes, j_addr_bytes...)

			interrupt_bytes = append(interrupt_bytes, current_interrupt_bytes...)

			current_interrupt_bytes = nil

		}

	}

	byte_index += uint32(len(interrupt_bytes) + int(PadSize))

	//Generate instruction bytecode

	instruction_bytes := []byte{}

	instruction_start_index := byte_index

	for _, v := range p.ProgramInstructions {

		instruction_bytes = append(instruction_bytes, generate_instruction_bytecode(
			v,
			data_block_addr,
			jump_block_addr,
		)...)

	}

	//Construct final byte array

	b_padding_bytes := []byte{PadValue, PadValue, PadValue, PadValue}

	b_data_block_start := make([]byte, 4)
	b_jmp_block_start := make([]byte, 4)
	b_int_block_start := make([]byte, 4)
	b_itn_block_start := make([]byte, 4)

	binary.LittleEndian.PutUint32(b_data_block_start, definition_start_index)
	binary.LittleEndian.PutUint32(b_jmp_block_start, jump_block_start_index)
	binary.LittleEndian.PutUint32(b_int_block_start, interrupt_block_start_index)
	binary.LittleEndian.PutUint32(b_itn_block_start, instruction_start_index)

	final_bytes = append(final_bytes, b_data_block_start...)
	final_bytes = append(final_bytes, b_jmp_block_start...)
	final_bytes = append(final_bytes, b_int_block_start...)
	final_bytes = append(final_bytes, b_itn_block_start...)

	final_bytes = append(final_bytes, b_padding_bytes...)

	final_bytes = append(final_bytes, definition_bytes...)
	final_bytes = append(final_bytes, b_padding_bytes...)

	final_bytes = append(final_bytes, jump_block_bytes...)
	final_bytes = append(final_bytes, b_padding_bytes...)

	final_bytes = append(final_bytes, interrupt_bytes...)
	final_bytes = append(final_bytes, b_padding_bytes...)

	final_bytes = append(final_bytes, instruction_bytes...)
	final_bytes = append(final_bytes, b_padding_bytes...)

	return final_bytes

}

//General purpose instruction for generating instruction bytecode

func generate_instruction_bytecode(i Instruction, d_block_addr map[string]uint32, j_blk_addr map[string]uint32) []byte {

	//TODO: sign bit
	//TODO: add offset for "hardware reserved" space

	var instruction_bytes []byte

	instruction_bytes = append(instruction_bytes,
		uint8(i.Instruction),
	)

	//Evaluate instruction args

	var addresses []uint32

	for _, v := range i.Data {
		var addr uint32

		if v[0] == '@' { //LDA or STA

			addr = d_block_addr[v[1:]]

		} else if i.Instruction == uint32(constants.IJump) || i.Instruction == uint32(constants.IConditionalJump) { //jump

			addr = uint32(j_blk_addr[v] + StackSize)

		} else if i.Instruction == uint32(constants.ICallInterrupt) {
			addr = uint32(constants.InterruptInts[v])
		} else {
			addr = constants.RegisterInts[v]
		}

		addresses = append(addresses, addr)
	}

	//Add args to byte array

	var data_array []byte

	if i.SingleData {

		data_array = make([]byte, 4)

		binary.LittleEndian.PutUint32(data_array[:], addresses[0])

	} else {
		data_array = make([]byte, 4)

		binary.LittleEndian.PutUint16(data_array[:], uint16(addresses[0]))
		binary.LittleEndian.PutUint16(data_array[2:], uint16(addresses[1]))
	}

	instruction_bytes = append(instruction_bytes, data_array...)

	return instruction_bytes

}
