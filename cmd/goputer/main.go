package main

// VM & Compiler CMD front end

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sccreeper/govm/pkg/compiler"
	"sccreeper/govm/pkg/constants"
	"sccreeper/govm/pkg/util"
	"strings"

	"github.com/fatih/color"
	"github.com/savioxavier/termlink"
	"github.com/urfave/cli/v2"
)

var use_json bool
var json_path string
var output_path string

var Commit string

var green_bold_underline = color.New([]color.Attribute{color.FgGreen, color.Bold, color.Underline}...)
var bold = color.New([]color.Attribute{color.Bold}...)
var underline = color.New([]color.Attribute{color.FgWhite, color.Underline}...)

func format_instruction(i_name string, i_data []string) string {

	return fmt.Sprintf("%s %s", color.GreenString(i_name), color.CyanString(strings.Join(i_data[:], " ")))

}

func convert_hex(i int) string {

	return fmt.Sprintf("0x"+"%08X", i)

}

func main() {

	app := &cli.App{
		Name:        "govmcd",
		Description: "Program that is the frontend for running VMs and compiling code",

		Commands: []*cli.Command{
			{
				Name:    "build",
				Aliases: []string{"b"},
				Usage:   "Used to compile programs",
				Action:  _compiler,

				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "json",
						Usage:       "Enable JSON outputting",
						Destination: &use_json,
					},
					&cli.StringFlag{
						Name:        "jsonpath",
						Usage:       "Output program structure/data in `FILE` ",
						Destination: &json_path,
					},
					&cli.StringFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						Usage:       "Output binary to `FILE`",
						Destination: &output_path,
					},
				},
			},
			{
				Name:    "disassemble",
				Aliases: []string{"d"},
				Usage:   "Used to disassemble programs",
				Action:  _disassemble,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func _compiler(ctx *cli.Context) error {

	file_path := ctx.Args().Get(0)

	// See if file exists

	if _, err := os.Stat(file_path); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	log.Printf("Compiling %s", file_path)

	//Read file

	data, err := os.ReadFile(file_path)
	util.CheckError(err)

	//Determine output path

	if ctx.String("output") == "" {
		output_path = strings.TrimSuffix(file_path, filepath.Ext(file_path))
	} else if ctx.String("output") != "" {
		output_path = ctx.String("output")
	}

	//Determine JSON path

	if use_json {
		if json_path == "" {
			json_path = fmt.Sprintf("%s.json", strings.TrimSuffix(output_path, filepath.Ext(output_path)))
		} else if use_json {
			json_path = ctx.String("jsonpath")
		}
	}

	log.Println(use_json)

	compiler_config := compiler.CompilerConfig{

		OutputPath: output_path,
		OutputJSON: use_json,
		JSONPath:   json_path,
	}

	err = compiler.Assemble(string(data), compiler_config)

	if err != nil {
		util.CheckError(err)
	}

	return nil
}

func _disassemble(ctx *cli.Context) error {

	file_path := ctx.Args().Get(0)

	data, err := os.ReadFile(file_path)
	util.CheckError(err)

	program, err := compiler.Disassemble(data, false)

	//Reverse map

	itn_map := make(map[constants.Instruction]string)

	for k, v := range constants.InstructionInts {

		itn_map[constants.Instruction(v)] = k

	}

	//Output disassembled program

	green_bold_underline.Printf("goputer Disassembler (%s)\n", termlink.Link(Commit[:10], fmt.Sprintf("https://github.com/sccreeper/goputer/commit/%s", Commit[0:10])))
	fmt.Println()

	underline.Println("Block addresses:")
	fmt.Println()

	color.White("Data block: %s", bold.Sprintf(convert_hex(int(program.StartIndexes[0]))))
	color.White("Jump blocks: %s", bold.Sprintf(convert_hex(int(program.StartIndexes[1]))))
	color.White("Interrupt table: %s", bold.Sprintf(convert_hex(int(program.StartIndexes[2]))))
	color.White("Instruction block: %s", bold.Sprintf(convert_hex(int(program.StartIndexes[3]))))

	fmt.Println()
	underline.Println("Definitions:")
	fmt.Println()

	defintion_byte_index := compiler.BlockAddrSize + compiler.PadSize

	for _, v := range program.ProgramDefinitions {

		// real_byte_index := ""

		fmt.Printf(
			"F: %s M: %s = %s\n",
			bold.Sprintf(convert_hex(int(defintion_byte_index))),
			bold.Sprintf(convert_hex(int(defintion_byte_index+compiler.StackSize))),
			strings.ReplaceAll(string(v), "\n", ""),
		)

		//Calculate memory address and index in file

		defintion_byte_index += uint32(len(v) + 4)

	}

	fmt.Println()
	underline.Println("Interrupt table")
	fmt.Println()

	//Output to console

	for k, v := range program.InterruptTable {

		fmt.Printf("%02d = %s\n", int(k), bold.Sprintf(convert_hex(int(v))))

	}

	fmt.Println()
	underline.Println("Jump blocks")
	fmt.Println()

	for k, v := range program.JumpBlocks {

		fmt.Println()
		bold.Printf("Jump %s (File: %s):\n", convert_hex(int(k+compiler.StackSize)), convert_hex(int(k)))
		fmt.Println()

		for _, v1 := range v {

			fmt.Println(format_instruction(itn_map[constants.Instruction(v1.Instruction)], v1.Data))

		}

	}

	fmt.Println()
	underline.Println("Instructions:")
	fmt.Println()

	for _, v := range program.Instructions {

		fmt.Println(format_instruction(itn_map[constants.Instruction(v.Instruction)], v.Data))

	}

	return err

}
