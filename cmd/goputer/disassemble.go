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

	file_path := ctx.Args().Get(0)

	data, err := os.ReadFile(file_path)
	util.CheckError(err)

	if string(data[:4]) != compiler.MagicString {
		fmt.Println("Error: Invalid file")
		os.Exit(1)
	}

	program, err := compiler.Disassemble(data, false)

	//Reverse map

	itn_map := make(map[constants.Instruction]string)

	for k, v := range constants.InstructionInts {

		itn_map[constants.Instruction(v)] = k

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

	defintion_byte_index := compiler.BlockAddrSize + compiler.PadSize

	for _, v := range program.ProgramDefinitions {

		// real_byte_index := ""

		fmt.Printf(
			"F: %s M: %s = %s\n",
			Bold.Sprintf(util.ConvertHex(int(defintion_byte_index))),
			Bold.Sprintf(util.ConvertHex(int(defintion_byte_index+compiler.StackSize))),
			strings.ReplaceAll(string(v), "\n", ""),
		)

		//Calculate memory address and index in file

		defintion_byte_index += uint32(len(v) + 4)

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

			fmt.Println(format_instruction(itn_map[constants.Instruction(v1.Instruction)], v1.Data))

		}

	}

	fmt.Println()
	Underline.Println("Instructions:")
	fmt.Println()

	for _, v := range program.Instructions {

		fmt.Println(format_instruction(itn_map[constants.Instruction(v.Instruction)], v.Data))

	}

	return err

}
