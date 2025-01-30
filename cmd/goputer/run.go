package main

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/util"

	"github.com/urfave/cli/v2"
)

// Run the program using a frontend
func _run(ctx *cli.Context) error {

	programBytes, err := os.ReadFile(GPExec)
	util.CheckError(err)

	if string(programBytes[:4]) != compiler.MagicString {
		fmt.Println("Error: Invalid file")
		os.Exit(1)
	}

	p, err := plugin.Open(fmt.Sprintf("./frontends/%s/%s%s", FrontendToUse, FrontendToUse, PluginExt))
	util.CheckError(err)

	runFunc, err := p.Lookup("Run")
	util.CheckError(err)

	runFunc.(func([]byte, []string))(programBytes, []string{filepath.Base(GPExec)})

	return nil

}
