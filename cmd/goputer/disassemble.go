package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"

	"github.com/fatih/color"
	"github.com/savioxavier/termlink"
	"github.com/urfave/cli/v2"
)

func _disassemble(ctx *cli.Context) error {

	filePath := ctx.Args().Get(0)

	data, err := os.ReadFile(filePath)
	util.CheckError(err)

	if string(data[:4]) != compiler.MagicString {
		fmt.Println("Error: Invalid file")
		os.Exit(1)
	}

	program, err := compiler.Disassemble(data, BeVerbose)

	//Reverse map

	itnMap := make(map[constants.Instruction]string)

	for k, v := range constants.InstructionInts {

		itnMap[constants.Instruction(v)] = k

	}

	//Output disassembled program

	GreenBoldUnderline.Printf("goputer Disassembler (%s)\n", termlink.Link(Commit[:10], fmt.Sprintf("https://github.com/sccreeper/goputer/commit/%s", Commit[0:10])))
	fmt.Println()

	Underline.Println("Block addresses:")
	fmt.Println()

	color.White("Interrupt table: %s", Bold.Sprintf(util.ConvertHex(int(program.StartIndexes[0]))))
	color.White("Data block: %s", Bold.Sprintf(util.ConvertHex(int(program.StartIndexes[1]))))
	color.White("Instructions: %s", Bold.Sprintf(util.ConvertHex(int(program.StartIndexes[2]))))
	color.White("Instruction entry point: %s", Bold.Sprintf(util.ConvertHex(int(program.StartIndexes[3]))))

	fmt.Println()
	Underline.Println("Definitions:")
	fmt.Println()

	definitionByteIndex := compiler.HeaderSize

	for _, v := range program.ProgramDefinitions {

		fmt.Printf(
			"File address: %s Memory address: %s Length: %s = %s\n",
			Bold.Sprintf(util.ConvertHex(int(definitionByteIndex))),
			Bold.Sprintf(util.ConvertHex(int(definitionByteIndex+compiler.StackSize))),
			Bold.Sprintf("%d", len(v)),
			hex.EncodeToString(v),
		)

		//Calculate memory address and index in file

		definitionByteIndex += uint32(len(v) + 4)

	}

	fmt.Println()
	Underline.Println("Interrupt table")
	fmt.Println()

	//Output to console

	for k, v := range program.InterruptTable {

		fmt.Printf("%02d = %s\n", int(k), Bold.Sprintf(util.ConvertHex(int(v))))

	}

	fmt.Println()
	Underline.Println("Instructions:")
	fmt.Println()

	for i, v := range program.Instructions {

		fmt.Printf(
			"%s: %s\n",
			Grey.Sprintf(util.ConvertHex((i*int(compiler.InstructionLength))+int(program.StartIndexes[2])+int(compiler.StackSize))),
			formatInstruction(
				itnMap[constants.Instruction(v.Instruction)],
				v.StringData,
			),
		)

	}

	return err

}
