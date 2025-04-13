package compiler

import (
	"encoding/binary"
	"sccreeper/goputer/pkg/constants"
	c "sccreeper/goputer/pkg/constants"
)

const (
	MagicString string = "GPTR"
)

// GenerateBytecode
//
// Takes in a ProgramStructure and returns the corresponding compiled bytecode.
func GenerateBytecode(p ProgramStructure) []byte {

	byteIndex := BlockAddrSize
	finalBytes := make([]byte, 0)

	byteIndex += uint32(len(constants.SubscribableInterrupts) * 6)

	definitionBlockAddresses := make(map[string]uint32)

	//Generate definition bytes first

	definitionBytes := []byte{}
	definitionStartIndex := byteIndex
	definitionAddrIndex := definitionStartIndex

	for i, d := range p.Definitions {

		definitionBlockAddresses[d.Name] = definitionAddrIndex + StackSize
		p.Definitions[i] = Definition{
			Name:       p.Definitions[i].Name,
			StringData: p.Definitions[i].StringData,
			ByteData:   p.Definitions[i].ByteData,
			Type:       p.Definitions[i].Type,

			Address: definitionAddrIndex + StackSize,
		}

		lengthBytes := make([]byte, 4)

		binary.LittleEndian.PutUint32(lengthBytes, uint32(len(d.ByteData)))

		dataBytes := d.ByteData

		definitionBytes = append(definitionBytes, lengthBytes...)
		definitionBytes = append(definitionBytes, dataBytes...)

		definitionAddrIndex += uint32(len(lengthBytes) + len(dataBytes))

	}

	//Increment the byte index

	byteIndex += uint32(len(definitionBytes))

	// Generate label addresses

	var labelAddresses map[string]uint32 = make(map[string]uint32)

	for k, v := range p.ProgramLabels {
		labelAddresses[k] = (uint32(v.InstructionOffset) * InstructionLength) + byteIndex + StackSize
	}

	// Generate interrupt jump table for existing interrupts

	interruptBytes := []byte{}

	interruptBlockStartIndex := byteIndex

	for i := c.IntMouseMove; i < c.IntKeyboardDown+1; i++ {

		var labelAddress uint32
		var interrupt c.Interrupt = c.Interrupt(i)

		if val, exists := p.InterruptSubscriptions[c.InterruptIntsReversed[i]]; exists {
			labelAddress = labelAddresses[val.LabelName]
		} else {
			labelAddress = 0
		}

		interruptTypeBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(interruptTypeBytes[:], uint16(interrupt))
		labelAdddressBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(labelAdddressBytes[:], uint32(labelAddress))

		interruptBytes = append(interruptBytes, interruptTypeBytes...)
		interruptBytes = append(interruptBytes, labelAdddressBytes...)

	}

	//Generate instruction bytecode

	instructionBytes := []byte{}

	for _, v := range p.ProgramInstructions {

		instructionBytes = append(instructionBytes, generateInstructionBytecode(
			v,
			definitionBlockAddresses,
			labelAddresses,
		)...)

	}

	//Construct final byte array

	finalBytes = append(finalBytes, []byte(MagicString)...)

	dataBlockStart := make([]byte, 4)
	interruptBlockStart := make([]byte, 4)
	instructionBlockStart := make([]byte, 4)

	binary.LittleEndian.PutUint32(dataBlockStart, definitionStartIndex)
	binary.LittleEndian.PutUint32(interruptBlockStart, interruptBlockStartIndex)
	binary.LittleEndian.PutUint32(instructionBlockStart, labelAddresses["start"])

	finalBytes = append(finalBytes, dataBlockStart...)
	finalBytes = append(finalBytes, interruptBlockStart...)
	finalBytes = append(finalBytes, instructionBlockStart...)

	finalBytes = append(finalBytes, interruptBytes...)

	finalBytes = append(finalBytes, definitionBytes...)

	finalBytes = append(finalBytes, instructionBytes...)
	finalBytes = append(finalBytes, []byte{0, 0, 0, 0}...)

	return finalBytes

}

// Generates individual instruction bytecode.
//
// 1 byte for instruction, 4 bytes for arguments.
func generateInstructionBytecode(i Instruction, definitionAddresses map[string]uint32, labelAddresses map[string]uint32) []byte {

	//TODO: sign bit
	//TODO: add offset for "hardware reserved" space

	var instructionBytes []byte

	instructionBytes = append(instructionBytes,
		uint8(i.Instruction),
	)

	//Evaluate instruction args

	var arguements []uint32

	for _, v := range i.StringData {
		var arg uint32

		if i.Instruction == uint32(c.IStore) || i.Instruction == uint32(c.ILoad) {

			if v[0] == '@' {
				arg = definitionAddresses[v[1:]]
			} else {
				arg = uint32(c.RegisterInts[v])
			}

		} else if i.Instruction == uint32(c.IJump) || i.Instruction == uint32(c.IConditionalJump) || i.Instruction == uint32(c.ICall) || i.Instruction == uint32(c.IConditionalCall) {

			if v[0] == '@' {
				arg = labelAddresses[v[1:]]
			} else {
				arg = uint32(c.RegisterInts[v])
			}

		} else if i.Instruction == uint32(c.ICallInterrupt) {
			arg = uint32(c.InterruptInts[v])
		} else {
			arg = c.RegisterInts[v]
		}

		arguements = append(arguements, arg)
	}

	//Add args to byte array

	var dataArray []byte

	if i.ArgumentCount == 0 {

		dataArray = []byte{0, 0, 0, 0}

	} else if i.ArgumentCount == 1 {

		dataArray = make([]byte, 4)

		binary.LittleEndian.PutUint32(dataArray[:], arguements[0])

	} else {
		dataArray = make([]byte, 4)

		binary.LittleEndian.PutUint16(dataArray[:], uint16(arguements[0]))
		binary.LittleEndian.PutUint16(dataArray[2:], uint16(arguements[1]))
	}

	instructionBytes = append(instructionBytes, dataArray...)

	return instructionBytes

}
