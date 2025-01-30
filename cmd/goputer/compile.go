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
var standaloneCodeGo string

type standaloneTemplate map[string]interface{}

// Main compile method for CLI
func _compiler(ctx *cli.Context) error {

	fmt.Printf("goputer compiler Version: %s\n", Commit[:10])

	filePath := ctx.Args().Get(0)

	// See if file exists

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	fmt.Printf("Compiling %s\n", filePath)

	prevDir, err := os.Getwd()
	util.CheckError(err)

	//Determine output path

	if ctx.String("output") == "" {
		OutputPath = strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".gp"
	} else if ctx.String("output") != "" {
		if filepath.Ext(ctx.String("output")) != ".gp" {
			OutputPath = ctx.String("output") + ".gp"
		} else {
			OutputPath = ctx.String("output")
		}
	}

	//Determine JSON path

	if UseJson {
		if JsonPath == "" {
			JsonPath = fmt.Sprintf("%s.json", strings.TrimSuffix(OutputPath, filepath.Ext(OutputPath)))
		} else if UseJson {
			JsonPath = ctx.String("jsonpath")
		}
	}

	compilerConfig := compiler.CompilerConfig{

		OutputPath: OutputPath,
		FilePath:   filepath.Base(filePath),
		OutputJSON: UseJson,
		JSONPath:   JsonPath,
		Verbose:    Verbose,
	}

	os.Chdir(filepath.Dir(filePath))

	//Assemble program & write to disk

	assembledProgram, err := compiler.Compile(compilerConfig.FilePath, getFile, compilerConfig, errorHandler)

	util.CheckError(err)

	os.Chdir(prevDir)

	//If standalone write to disk differently
	if IsStandalone {

		standaloneBytes := _standalone(assembledProgram.ProgramBytes)
		tempFile, err := os.Create("./alone_temp.go")
		util.CheckError(err)

		tempFile.Write(standaloneBytes)

		os.Remove(compilerConfig.OutputPath)

		//Calculate checksums before compression

		programHash := sha256.New()
		programHash.Write(assembledProgram.ProgramBytes)

		pluginFile, err := os.ReadFile(fmt.Sprintf("./frontends/%s/%s%s", FrontendToUse, FrontendToUse, PluginExt))
		util.CheckError(err)

		pluginHash := sha256.New()
		pluginHash.Write(pluginFile)

		//Compress files
		var b bytes.Buffer

		pluginCompressed, err := os.Create("plugin_compressed")
		util.CheckError(err)

		w := zlib.NewWriter(&b)
		w.Write(pluginFile)
		w.Close()
		pluginCompressed.Write(b.Bytes())
		b.Reset()

		codeCompressed, err := os.Create("code_compressed")
		util.CheckError(err)

		w = zlib.NewWriter(&b)
		w.Write(assembledProgram.ProgramBytes)
		w.Close()
		codeCompressed.Write(b.Bytes())
		b.Reset()

		ldFlags := fmt.Sprintf(
			"-s -w -X main.ProgramCheck=%s -X main.PluginCheck=%s",
			hex.EncodeToString(programHash.Sum(nil)),
			hex.EncodeToString(pluginHash.Sum(nil)),
		)

		cmd := exec.Command("go", "build", "-ldflags", ldFlags, "-o", compilerConfig.OutputPath, "./alone_temp.go")
		var out bytes.Buffer

		cmd.Stdout = &out

		cmd.Run()

		//Cleanup
		//os.Remove("alone_temp.go")
		//os.Remove("alone_program_bytes")

		log.Printf("Finished making executable %s", compilerConfig.OutputPath)

	} else {
		os.WriteFile(compilerConfig.OutputPath, assembledProgram.ProgramBytes, 0666)
	}

	//JSON

	if compilerConfig.OutputJSON {
		err = os.WriteFile(JsonPath, []byte(assembledProgram.ProgramJson), 0666)

		util.CheckError(err)

		log.Printf("Outputted JSON structure to '%s'", JsonPath)
	}

	return nil
}

// Extra method for standalone executables.
func _standalone(program []byte) []byte {

	log.Println("Making standalone executable...")

	var finalCode bytes.Buffer

	bytesFile, err := os.Create("alone_program_bytes")
	util.CheckError(err)
	bytesFile.Write(program)

	t := template.Must(template.New("").Parse(standaloneCodeGo))
	t.Execute(
		&finalCode,
		standaloneTemplate{
			"plugin":       "plugin_compressed",
			"program_code": "code_compressed",
		},
	)

	return finalCode.Bytes()
}

func errorHandler(errorType compiler.ErrorType, errorText string) {

	fmt.Println(errorText)
	os.Exit(1)

}

func getFile(path string) []byte {

	f, err := os.ReadFile(path)

	util.CheckError(err)

	return f

}
