package compiler

import (
	"encoding/binary"
	"log"
	c "sccreeper/goputer/pkg/constants"
	"strconv"
)

const (
	MagicString string = "GPTR"
)

// GenerateBytecode
//
// Takes in a ProgramStructure and returns the corresponding compiled bytecode.
func GenerateBytecode(p ProgramStructure, verbose bool) []byte {

	byteIndex := HeaderSize
	finalBytes := make([]byte, 0)

	byteIndex += uint32(len(c.SubscribableInterrupts) * 6)

	definitionBlockAddresses := make(map[string]uint32)

	//Generate definition bytes first

	definitionBytes := make([]byte, 0)
	definitionStartIndex := byteIndex
	definitionAddrIndex := definitionStartIndex

	if verbose {
		log.Printf("Definition start address is %d", definitionStartIndex)
	}

	for i, d := range p.Definitions {

		definitionBlockAddresses[d.Name] = definitionAddrIndex + MemOffset
		p.Definitions[i] = Definition{
			Name:       p.Definitions[i].Name,
			StringData: p.Definitions[i].StringData,
			ByteData:   p.Definitions[i].ByteData,
			Type:       p.Definitions[i].Type,

			Address: definitionAddrIndex + MemOffset,
		}

		lengthBytes := make([]byte, 4)

		binary.LittleEndian.PutUint32(lengthBytes, uint32(len(d.ByteData)))

		definitionBytes = append(definitionBytes, lengthBytes...)
		definitionBytes = append(definitionBytes, d.ByteData...)

		definitionAddrIndex += uint32(len(lengthBytes) + len(d.ByteData))

	}

	//Increment the byte index

	byteIndex += uint32(len(definitionBytes))

	// Generate label addresses

	var labelAddresses map[string]uint32 = make(map[string]uint32)

	for k, v := range p.ProgramLabels {
		labelAddresses[k] = (uint32(v.InstructionOffset) * InstructionLength) + byteIndex + MemOffset
	}

	// Generate interrupt jump table for all interrupts

	interruptBytes := []byte{}

	interruptBlockStartIndex := HeaderSize

	for _, v := range c.SubscribableInterrupts {

		var labelAddress uint32
		var interrupt c.Interrupt = v

		if val, exists := p.InterruptSubscriptions[c.InterruptIntsReversed[v]]; exists {
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

	if verbose {
		log.Printf("Interrupt table is %d bytes long", len(interruptBytes))
	}

	//Generate instruction bytecode

	instructionBytes := []byte{}

	for _, v := range p.ProgramInstructions {

		instructionBytes = append(instructionBytes, generateInstructionBytecode(
			v,
			definitionBlockAddresses,
			labelAddresses,
			byteIndex + MemOffset,
		)...)

	}

	//Construct final byte array

	finalBytes = append(finalBytes, []byte(MagicString)...)

	definitionBlockStart := make([]byte, 4)
	interruptBlockStart := make([]byte, 4)
	instructionBlockStart := make([]byte, 4)
	instructionEntryPoint := make([]byte, 4)

	binary.LittleEndian.PutUint32(definitionBlockStart, definitionStartIndex)
	binary.LittleEndian.PutUint32(interruptBlockStart, interruptBlockStartIndex)
	binary.LittleEndian.PutUint32(instructionBlockStart, definitionStartIndex+uint32(len(definitionBytes)))
	binary.LittleEndian.PutUint32(instructionEntryPoint, labelAddresses["start"]-MemOffset)

	finalBytes = append(finalBytes, interruptBlockStart...)
	finalBytes = append(finalBytes, definitionBlockStart...)
	finalBytes = append(finalBytes, instructionBlockStart...)
	finalBytes = append(finalBytes, instructionEntryPoint...)

	finalBytes = append(finalBytes, interruptBytes...)

	finalBytes = append(finalBytes, definitionBytes...)

	finalBytes = append(finalBytes, instructionBytes...)
	finalBytes = append(finalBytes, []byte{0, 0, 0, 0}...)

	return finalBytes

}

// Generates individual instruction bytecode.
//
// 1 byte for instruction, 4 bytes for arguments.
func generateInstructionBytecode(itn Instruction, definitionAddresses map[string]uint32, labelAddresses map[string]uint32, memOffset uint32) []byte {

	//TODO: sign bit
	//TODO: add offset for "hardware reserved" space

	var instructionBytes []byte

	instructionBytes = append(instructionBytes,
		uint8(itn.Instruction),
	)

	//Evaluate instruction args

	var arguments []uint32

	for _, stringArg := range itn.StringData {
		var arg uint32

		// Do immediate val ahead of time, if there is one
		if stringArg[0] == '$' {
			if stringArg[1] == ':' {
				
				x, _ := strconv.Atoi(stringArg[2:])
				arg = uint32(x)

				arg += memOffset

			} else {
				x, _ := strconv.Atoi(stringArg[1:])
				arg = uint32(x)
			}

			arguments = append(arguments, arg)

			continue
		}

		if itn.Instruction == uint32(c.IStore) || itn.Instruction == uint32(c.ILoad) {

			if stringArg[0] == '@' {
				arg = definitionAddresses[stringArg[1:]]
			} else {
				arg = uint32(c.RegisterInts[stringArg])
			}

		} else if itn.Instruction == uint32(c.IJump) || itn.Instruction == uint32(c.IConditionalJump) || itn.Instruction == uint32(c.ICall) || itn.Instruction == uint32(c.IConditionalCall) {

			if stringArg[0] == '@' {
				arg = labelAddresses[stringArg[1:]]
			} else {
				arg = uint32(c.RegisterInts[stringArg])
			}

		} else if itn.Instruction == uint32(c.ICallInterrupt) {
			arg = uint32(c.InterruptInts[stringArg])
		} else {

			arg = c.RegisterInts[stringArg]
		
		}

		arguments = append(arguments, arg)
	}

	if itn.HasImmediate {

		if len(arguments) == 1 {
			
			instructionBytes = binary.LittleEndian.AppendUint32(instructionBytes, arguments[0])
		
		} else {
			if itn.ImmediateIndex == 0 {
				instructionBytes[0] |= byte(c.ItnFlagFirstArgImmediate)
			} else if itn.ImmediateIndex == 1 {
				instructionBytes[0] |= byte(c.ItnFlagSecondArgImmediate)
			}

			// Process immediate

			var argValue uint32

			argValue = arguments[itn.ImmediateIndex]

			if itn.ImmediateIndex == 0 {
				argValue |= arguments[1] << 26	
			} else {
				argValue |= arguments[0] << 26	
			}

			instructionBytes = binary.LittleEndian.AppendUint32(instructionBytes, argValue)
		}

		return instructionBytes
	
	} else {
		//Add args to byte array

		var dataArray []byte

		if itn.ArgumentCount == 0 {

			dataArray = []byte{0, 0, 0, 0}

		} else if itn.ArgumentCount == 1 {

			dataArray = make([]byte, 4)

			binary.LittleEndian.PutUint32(dataArray[:], arguments[0])

		} else {
			dataArray = make([]byte, 4)

			binary.LittleEndian.PutUint16(dataArray[:], uint16(arguments[0]))
			binary.LittleEndian.PutUint16(dataArray[2:], uint16(arguments[1]))
		}

		instructionBytes = append(instructionBytes, dataArray...)

		return instructionBytes	
	}

}
