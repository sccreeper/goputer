package main

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sccreeper/goputer/pkg/util"

	"github.com/urfave/cli/v2"
)

// Run the program using a frontend
func _run(ctx *cli.Context) error {

	program_bytes, err := os.ReadFile(GPExec)
	util.CheckError(err)

	p, err := plugin.Open(fmt.Sprintf("./frontends/%s/%s%s", FrontendToUse, FrontendToUse, PluginExt))
	util.CheckError(err)

	run_func, err := p.Lookup("Run")
	util.CheckError(err)

	run_func.(func([]byte, []string))(program_bytes, []string{filepath.Base(GPExec)})

	return nil

}
