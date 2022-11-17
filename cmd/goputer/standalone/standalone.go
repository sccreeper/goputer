package main

import (
	_ "embed"
	"os"
	"plugin"
)

//go:embed {{.plugin}}
var plugin_bytes []byte

//go:embed {{.program_code}}
var program_code []byte

func main() {
	plugin_file, err := os.CreateTemp(os.TempDir(), "*.so")
	check_err(err)

	plugin_file.Write(plugin_bytes)

	p, err := plugin.Open(plugin_file.Name())
	check_err(err)

	run_func, err := p.Lookup("Run")

	run_func.(func([]byte, []string))(program_code, []string{})

}

func check_err(err error) {

	if err != nil {
		panic(err)
	}

}
