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

	"github.com/urfave/cli/v2"
)

var use_json bool
var json_path string
var output_path string

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

	//Reverse dict

	itn_dict := make(map[constants.Instruction]string)

	for k, v := range constants.InstructionInts {

		itn_dict[constants.Instruction(v)] = k

	}

	//Output disassembled program

	log.Println("Disassembled program structure")

	log.Println("Block addresses:")

	log.Printf("Data block: %s", convert_hex(int(program.StartIndexes[0])))
	log.Printf("Jump blocks: %s", convert_hex(int(program.StartIndexes[1])))
	log.Printf("Interrupt table: %s", convert_hex(int(program.StartIndexes[2])))
	log.Printf("Instruction block: %s", convert_hex(int(program.StartIndexes[3])))

	log.Println("Definitions:")

	for index, v := range program.ProgramDefinitions {

		log.Printf("Def %s = %s", convert_hex(index), v)

	}

	log.Println("Interrupt table")

	//Output to console

	for k, v := range program.InterruptTable {

		log.Printf("Int %d = %s", int(k), convert_hex(int(v)))

	}

	log.Println("Jump blocks")

	for k, v := range program.JumpBlocks {

		log.Printf("Jump %s:", convert_hex(int(k)))

		for _, v1 := range v {

			log.Printf("%s %s", itn_dict[constants.Instruction(v1.Instruction)], strings.Join(v1.Data[:], " "))

		}

		log.Println("============================")

	}

	log.Println("Instructions:")

	for _, v := range program.Instructions {

		log.Printf("%s %s", itn_dict[constants.Instruction(v.Instruction)], strings.Join(v.Data[:], " "))

	}

	return err

}
