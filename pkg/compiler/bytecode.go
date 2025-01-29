package compiler

import (
	"encoding/binary"
	"sccreeper/goputer/pkg/constants"

	"golang.org/x/exp/slices"
)

const (
	MagicString string = "GPTR"
)

// GenerateBytecode
//
// Takes in a ProgramStructure and returns the corresponding compiled bytecode.
func GenerateBytecode(p ProgramStructure) []byte {

	byteIndex := BlockAddrSize + 4
	finalBytes := make([]byte, 0)

	jumpBlockAddr := make(map[string]uint32)
	dataBlockAddr := make(map[string]uint32)

	//Generate definition bytecode first

	definitionBytes := []byte{}
	definitionStartIndex := byteIndex
	definitionAddrIndex := definitionStartIndex

	for _, d := range p.Definitions {

		dataBlockAddr[d.Name] = definitionAddrIndex + StackSize

		lengthBytes := make([]byte, 4)

		binary.LittleEndian.PutUint32(lengthBytes, uint32(len(d.Data)))

		dataBytes := d.Data

		definitionBytes = append(definitionBytes, lengthBytes...)
		definitionBytes = append(definitionBytes, dataBytes...)

		definitionAddrIndex += uint32(len(lengthBytes) + len(dataBytes))

	}

	//Increment the byte index

	byteIndex += uint32(len(definitionBytes)) + PadSize

	//Generate the jump block block

	jumpBlockStartIndex := byteIndex
	jumpBlockAddrIndex := jumpBlockStartIndex

	jumpBlockBytes := []byte{}

	//Order keys first

	for _, s := range p.InstructionBlockNames {

		v := p.InstructionBlocks[s]

		currentJumpBlockBytes := []byte{}

		jumpBlockAddr[v.Name] = jumpBlockAddrIndex

		for _, i := range v.Instructions {

			currentJumpBlockBytes = append(currentJumpBlockBytes, generateInstructionBytecode(i, dataBlockAddr, jumpBlockAddr)...)

		}

		currentJumpBlockBytes = append(currentJumpBlockBytes, []byte{0, 0, 0, 0, 0}...)

		jumpBlockBytes = append(jumpBlockBytes, currentJumpBlockBytes...)

		jumpBlockAddrIndex += uint32(len(currentJumpBlockBytes))

		currentJumpBlockBytes = nil

	}

	byteIndex += uint32(len(jumpBlockBytes)) + PadSize

	//Generate the interrupt table

	includedInterrupts := []constants.Interrupt{}
	currentInterruptBytes := []byte{}
	interruptBytes := []byte{}

	interruptBlockStartIndex := byteIndex

	for _, v := range p.InterruptSubscriptions {

		i := v.Interrupt
		jAddr := jumpBlockAddr[v.JumpBlockName]

		iBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(iBytes[:], uint16(i))
		jAddrBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(jAddrBytes[:], jAddr+StackSize)

		currentInterruptBytes = append(currentInterruptBytes, iBytes...)
		currentInterruptBytes = append(currentInterruptBytes, jAddrBytes...)

		interruptBytes = append(interruptBytes, currentInterruptBytes...)

		currentInterruptBytes = nil

		includedInterrupts = append(includedInterrupts, v.Interrupt)

	}

	//Generate the rest of the interrupt table bytes

	for _, v := range constants.SubscribableInterrupts {

		if slices.Contains(includedInterrupts, v) {
			continue
		} else {
			i := v
			jAddr := uint32(0x00000000)

			iBytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(iBytes[:], uint16(i))
			jAddrBytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(jAddrBytes[:], jAddr)

			currentInterruptBytes = append(currentInterruptBytes, iBytes...)
			currentInterruptBytes = append(currentInterruptBytes, jAddrBytes...)

			interruptBytes = append(interruptBytes, currentInterruptBytes...)

			currentInterruptBytes = nil

		}

	}

	byteIndex += uint32(len(interruptBytes) + int(PadSize))

	//Generate instruction bytecode

	instructionBytes := []byte{}

	instructionStartIndex := byteIndex

	for _, v := range p.ProgramInstructions {

		instructionBytes = append(instructionBytes, generateInstructionBytecode(
			v,
			dataBlockAddr,
			jumpBlockAddr,
		)...)

	}

	//Construct final byte array

	finalBytes = append(finalBytes, []byte(MagicString)...)

	bPaddingBytes := []byte{PadValue, PadValue, PadValue, PadValue}

	bDataBlockStart := make([]byte, 4)
	bJmpBlockStart := make([]byte, 4)
	bIntBlockStart := make([]byte, 4)
	bItnBlockStart := make([]byte, 4)

	binary.LittleEndian.PutUint32(bDataBlockStart, definitionStartIndex)
	binary.LittleEndian.PutUint32(bJmpBlockStart, jumpBlockStartIndex)
	binary.LittleEndian.PutUint32(bIntBlockStart, interruptBlockStartIndex)
	binary.LittleEndian.PutUint32(bItnBlockStart, instructionStartIndex)

	finalBytes = append(finalBytes, bDataBlockStart...)
	finalBytes = append(finalBytes, bJmpBlockStart...)
	finalBytes = append(finalBytes, bIntBlockStart...)
	finalBytes = append(finalBytes, bItnBlockStart...)

	finalBytes = append(finalBytes, bPaddingBytes...)

	finalBytes = append(finalBytes, definitionBytes...)
	finalBytes = append(finalBytes, bPaddingBytes...)

	finalBytes = append(finalBytes, jumpBlockBytes...)
	finalBytes = append(finalBytes, bPaddingBytes...)

	finalBytes = append(finalBytes, interruptBytes...)
	finalBytes = append(finalBytes, bPaddingBytes...)

	finalBytes = append(finalBytes, instructionBytes...)
	finalBytes = append(finalBytes, bPaddingBytes...)

	return finalBytes

}

// Generates individual instruction bytecode.
//
// 1 byte for instruction, 4 bytes for arguments.
func generateInstructionBytecode(i Instruction, dBlockAddr map[string]uint32, jBlkAddr map[string]uint32) []byte {

	//TODO: sign bit
	//TODO: add offset for "hardware reserved" space

	var instructionBytes []byte

	instructionBytes = append(instructionBytes,
		uint8(i.Instruction),
	)

	//Evaluate instruction args

	var addresses []uint32

	for _, v := range i.Data {
		var addr uint32

		if v[0] == '@' { //LDA or STA

			addr = dBlockAddr[v[1:]]

		} else if i.Instruction == uint32(constants.IJump) || i.Instruction == uint32(constants.IConditionalJump) || i.Instruction == uint32(constants.ICall) || i.Instruction == uint32(constants.IConditionalCall) { //jump

			addr = uint32(jBlkAddr[v] + StackSize)

		} else if i.Instruction == uint32(constants.ICallInterrupt) {
			addr = uint32(constants.InterruptInts[v])
		} else {
			addr = constants.RegisterInts[v]
		}

		addresses = append(addresses, addr)
	}

	//Add args to byte array

	var dataArray []byte

	if i.SingleData {

		dataArray = make([]byte, 4)

		binary.LittleEndian.PutUint32(dataArray[:], addresses[0])

	} else {
		dataArray = make([]byte, 4)

		binary.LittleEndian.PutUint16(dataArray[:], uint16(addresses[0]))
		binary.LittleEndian.PutUint16(dataArray[2:], uint16(addresses[1]))
	}

	instructionBytes = append(instructionBytes, dataArray...)

	return instructionBytes

}
