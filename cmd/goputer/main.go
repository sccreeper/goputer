package main

// VM & Compiler CMD front end

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"runtime"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/savioxavier/termlink"
	"github.com/urfave/cli/v2"

	_ "embed"
)

//go:embed standalone/standalone.go
var standalone_code_go string

type standalone_template map[string]interface{}

var use_json bool
var json_path string
var output_path string
var verbose bool = false

var frontend_to_use string
var gp_exec string
var is_standalone bool

var Commit string

var green_bold_underline = color.New([]color.Attribute{color.FgGreen, color.Bold, color.Underline}...)
var bold = color.New([]color.Attribute{color.Bold}...)
var underline = color.New([]color.Attribute{color.FgWhite, color.Underline}...)

var plugin_ext string

func format_instruction(i_name string, i_data []string) string {

	return fmt.Sprintf("%s %s", color.GreenString(i_name), color.CyanString(strings.Join(i_data[:], " ")))

}

func main() {

	if runtime.GOOS == "windows" {
		plugin_ext = ".dll"
	} else {
		plugin_ext = ".so"
	}

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
					&cli.BoolFlag{
						Name:        "verbose",
						Aliases:     []string{"v"},
						Usage:       "Verbose log output",
						Destination: &verbose,
					},
					&cli.BoolFlag{
						Name:        "standalone",
						Aliases:     []string{"s"},
						Usage:       "Create a standalone executable",
						Destination: &is_standalone,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "frontend",
						Aliases:     []string{"f"},
						Usage:       "Frontend to create standalone with",
						Destination: &frontend_to_use,
						Required:    false,
					},
				},
			},
			{
				Name:    "disassemble",
				Aliases: []string{"d"},
				Usage:   "Used to disassemble programs",
				Action:  _disassemble,
			},
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "Run programs",
				Action:  _run,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "frontend",
						Aliases: []string{"f"},
						Usage:   "Frontend to use",
						// DefaultText: "gp32",
						Destination: &frontend_to_use,
					},
					&cli.StringFlag{
						Name:        "exec",
						Aliases:     []string{"e"},
						Usage:       "Executable to run",
						Destination: &gp_exec,
					},
				},
			},
			{
				Name:   "list",
				Usage:  "Lists plugins available",
				Action: _list_plugins,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func _compiler(ctx *cli.Context) error {

	log.Printf("goputer compiler Version: %s", Commit[:10])

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

	compiler_config := compiler.CompilerConfig{

		OutputPath: output_path,
		OutputJSON: use_json,
		JSONPath:   json_path,
		Verbose:    verbose,
	}

	//Assemble program & write to disk

	assembled_program, err := compiler.Assemble(string(data), compiler_config)
	util.CheckError(err)

	//If standlone write to disk differently
	if is_standalone {

		standalone_bytes := _standalone(assembled_program.ProgramBytes)
		temp_file, err := os.Create("./alone_temp.go")
		util.CheckError(err)

		temp_file.Write(standalone_bytes)

		os.Remove(compiler_config.OutputPath)

		//Calculate checksums before compression

		program_hash := sha256.New()
		program_hash.Write(assembled_program.ProgramBytes)

		plugin_file, err := os.ReadFile(fmt.Sprintf("./frontends/%s/%s%s", frontend_to_use, frontend_to_use, plugin_ext))
		util.CheckError(err)

		plugin_hash := sha256.New()
		plugin_hash.Write(plugin_file)

		//Compress files
		var b bytes.Buffer

		plugin_compressed, err := os.Create("plugin_compressed")
		util.CheckError(err)

		w := zlib.NewWriter(&b)
		w.Write(plugin_file)
		w.Close()
		plugin_compressed.Write(b.Bytes())
		b.Reset()

		code_compressed, err := os.Create("code_compressed")
		util.CheckError(err)

		w = zlib.NewWriter(&b)
		w.Write(assembled_program.ProgramBytes)
		w.Close()
		code_compressed.Write(b.Bytes())
		b.Reset()

		ld_flags := fmt.Sprintf(
			"-s -w -X main.ProgramCheck=%s -X main.PluginCheck=%s",
			hex.EncodeToString(program_hash.Sum(nil)),
			hex.EncodeToString(plugin_hash.Sum(nil)),
		)

		cmd := exec.Command("go", "build", "-ldflags", ld_flags, "-o", compiler_config.OutputPath, "./alone_temp.go")
		var out bytes.Buffer

		cmd.Stdout = &out

		cmd.Run()

		//Cleanup
		//os.Remove("alone_temp.go")
		//os.Remove("alone_program_bytes")

		log.Printf("Finished making executable %s", compiler_config.OutputPath)

	} else {
		os.WriteFile(compiler_config.OutputPath, assembled_program.ProgramBytes, 0666)
	}

	//JSON

	if compiler_config.OutputJSON {
		err = os.WriteFile(json_path, []byte(assembled_program.ProgramJson), 0666)

		util.CheckError(err)

		log.Printf("Outputted JSON structure to '%s'", json_path)
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

	color.White("Data block: %s", bold.Sprintf(util.ConvertHex(int(program.StartIndexes[0]))))
	color.White("Jump blocks: %s", bold.Sprintf(util.ConvertHex(int(program.StartIndexes[1]))))
	color.White("Interrupt table: %s", bold.Sprintf(util.ConvertHex(int(program.StartIndexes[2]))))
	color.White("Instruction block: %s", bold.Sprintf(util.ConvertHex(int(program.StartIndexes[3]))))

	fmt.Println()
	underline.Println("Definitions:")
	fmt.Println()

	defintion_byte_index := compiler.BlockAddrSize + compiler.PadSize

	for _, v := range program.ProgramDefinitions {

		// real_byte_index := ""

		fmt.Printf(
			"F: %s M: %s = %s\n",
			bold.Sprintf(util.ConvertHex(int(defintion_byte_index))),
			bold.Sprintf(util.ConvertHex(int(defintion_byte_index+compiler.StackSize))),
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

		fmt.Printf("%02d = %s\n", int(k), bold.Sprintf(util.ConvertHex(int(v))))

	}

	fmt.Println()
	underline.Println("Jump blocks")
	fmt.Println()

	for k, v := range program.JumpBlocks {

		fmt.Println()
		bold.Printf("Jump %s (File: %s):\n", util.ConvertHex(int(k+compiler.StackSize)), util.ConvertHex(int(k)))
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

// Run the program using the default frontend
func _run(ctx *cli.Context) error {

	program_bytes, err := os.ReadFile(gp_exec)
	util.CheckError(err)

	p, err := plugin.Open(fmt.Sprintf("./frontends/%s/%s%s", frontend_to_use, frontend_to_use, plugin_ext))
	util.CheckError(err)

	run_func, err := p.Lookup("Run")
	util.CheckError(err)

	run_func.(func([]byte, []string))(program_bytes, []string{filepath.Base(gp_exec)})

	return nil

}

func _list_plugins(ctx *cli.Context) error {

	plugin_dir, err := ioutil.ReadDir("./frontends/")
	util.CheckError(err)

	for _, v := range plugin_dir {

		p, err := plugin.Open(fmt.Sprintf("./frontends/%s/%s%s", v.Name(), v.Name(), plugin_ext))
		util.CheckError(err)

		_name, err := p.Lookup("Name")
		util.CheckError(err)
		description, err := p.Lookup("Description")
		util.CheckError(err)
		authour, err := p.Lookup("Authour")
		util.CheckError(err)
		repo, err := p.Lookup("Repository")
		util.CheckError(err)

		fmt.Println()
		bold.Print(*_name.(*string) + "\n")
		fmt.Println()

		fmt.Printf("%s %s\n", bold.Sprintf("Description:"), *description.(*string))
		fmt.Printf("%s %s\n", bold.Sprintf("Authour:"), *authour.(*string))
		fmt.Printf("%s %s\n", bold.Sprintf("Repository:"), *repo.(*string))

	}

	fmt.Println()

	fmt.Printf("Found %d frontend(s)", len(plugin_dir))

	return nil
}

func _standalone(program []byte) []byte {

	log.Println("Making standalone executable...")

	var final_code bytes.Buffer

	bytes_file, err := os.Create("alone_program_bytes")
	util.CheckError(err)
	bytes_file.Write(program)

	t := template.Must(template.New("").Parse(standalone_code_go))
	t.Execute(
		&final_code,
		standalone_template{
			"plugin":       "plugin_compressed",
			"program_code": "code_compressed",
		},
	)

	return final_code.Bytes()
}
