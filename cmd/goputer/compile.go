package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/util"
	"strings"

	"github.com/urfave/cli/v2"

	_ "embed"
)

//Standalone variables.

//go:embed standalone/standalone.go
var standalone_code_go string

type standalone_template map[string]interface{}

// Main compile method for CLI
func _compiler(ctx *cli.Context) error {

	fmt.Printf("goputer compiler Version: %s\n", Commit[:10])

	file_path := ctx.Args().Get(0)

	// See if file exists

	if _, err := os.Stat(file_path); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	fmt.Printf("Compiling %s\n", file_path)

	//Read file

	data, err := os.ReadFile(file_path)
	util.CheckError(err)

	//Determine output path

	if ctx.String("output") == "" {
		OutputPath = strings.TrimSuffix(file_path, filepath.Ext(file_path))
	} else if ctx.String("output") != "" {
		OutputPath = ctx.String("output")
	}

	//Determine JSON path

	if UseJson {
		if JsonPath == "" {
			JsonPath = fmt.Sprintf("%s.json", strings.TrimSuffix(OutputPath, filepath.Ext(OutputPath)))
		} else if UseJson {
			JsonPath = ctx.String("jsonpath")
		}
	}

	compiler_config := compiler.CompilerConfig{

		OutputPath: OutputPath,
		FilePath:   file_path,
		OutputJSON: UseJson,
		JSONPath:   JsonPath,
		Verbose:    Verbose,
	}

	//Assemble program & write to disk

	assembled_program, err := compiler.Compile(string(data), compiler_config, error_handler)
	util.CheckError(err)

	//If standlone write to disk differently
	if IsStandalone {

		standalone_bytes := _standalone(assembled_program.ProgramBytes)
		temp_file, err := os.Create("./alone_temp.go")
		util.CheckError(err)

		temp_file.Write(standalone_bytes)

		os.Remove(compiler_config.OutputPath)

		//Calculate checksums before compression

		program_hash := sha256.New()
		program_hash.Write(assembled_program.ProgramBytes)

		plugin_file, err := os.ReadFile(fmt.Sprintf("./frontends/%s/%s%s", FrontendToUse, FrontendToUse, PluginExt))
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
		err = os.WriteFile(JsonPath, []byte(assembled_program.ProgramJson), 0666)

		util.CheckError(err)

		log.Printf("Outputted JSON structure to '%s'", JsonPath)
	}

	return nil
}

// Extra method for standalone executables.
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

func error_handler(error_type compiler.ErrorType, error_text string) {

	fmt.Println(error_text)
	os.Exit(1)

}
