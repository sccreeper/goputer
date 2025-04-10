package main

import (
	"fmt"
	"os"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"
	"strings"

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

	program, err := compiler.Disassemble(data, false)

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

	color.White("Data block: %s", Bold.Sprintf(util.ConvertHex(int(program.StartIndexes[0]))))
	color.White("Jump blocks: %s", Bold.Sprintf(util.ConvertHex(int(program.StartIndexes[1]))))
	color.White("Interrupt table: %s", Bold.Sprintf(util.ConvertHex(int(program.StartIndexes[2]))))
	color.White("Instruction block: %s", Bold.Sprintf(util.ConvertHex(int(program.StartIndexes[3]))))

	fmt.Println()
	Underline.Println("Definitions:")
	fmt.Println()

	definitionByteIndex := compiler.BlockAddrSize + compiler.PadSize

	for _, v := range program.ProgramDefinitions {

		fmt.Printf(
			"F: %s M: %s = %s\n",
			Bold.Sprintf(util.ConvertHex(int(definitionByteIndex))),
			Bold.Sprintf(util.ConvertHex(int(definitionByteIndex+compiler.StackSize))),
			strings.ReplaceAll(string(v), "\n", ""),
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
	Underline.Println("Jump blocks")
	fmt.Println()

	for k, v := range program.JumpBlocks {

		fmt.Println()
		Bold.Printf("Jump %s (File: %s):\n", util.ConvertHex(int(k+compiler.StackSize)), util.ConvertHex(int(k)))
		fmt.Println()

		for _, v1 := range v {

			fmt.Println(
				formatInstruction(
					itnMap[constants.Instruction(v1.Instruction)],
					v1.Data,
				),
			)

		}

	}

	fmt.Println()
	Underline.Println("Instructions:")
	fmt.Println()

	for _, v := range program.Instructions {

		fmt.Println(
			formatInstruction(
				itnMap[constants.Instruction(v.Instruction)],
				v.Data,
			),
		)

	}

	return err

}
