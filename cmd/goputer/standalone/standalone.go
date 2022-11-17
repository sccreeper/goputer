package main

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"log"
	"os"
	"plugin"
)

//go:embed {{.plugin}}
var plugin_bytes []byte

//go:embed {{.program_code}}
var program_code []byte

var ProgramCheck string
var PluginCheck string

func main() {
	plugin_file, err := os.CreateTemp(os.TempDir(), "*.so")
	check_err(err)

	plugin_file.Write(plugin_bytes)

	//Verify file integrity

	if ProgramCheck != checksum(program_code) {
		log.Println("Invalid checksum!")
		log.Println(checksum(program_code))
		os.Exit(1)
	}

	p, err := plugin.Open(plugin_file.Name())
	check_err(err)

	b, err := os.ReadFile(plugin_file.Name())

	if PluginCheck != checksum(b) || PluginCheck != checksum(plugin_bytes) {
		log.Println("Invalid checksum")
		os.Exit(1)
	}

	run_func, err := p.Lookup("Run")

	run_func.(func([]byte, []string))(program_code, []string{})

}

func check_err(err error) {

	if err != nil {
		panic(err)
	}

}

func checksum(b []byte) string {

	h := sha256.New()
	h.Write(b)

	return hex.EncodeToString(h.Sum(nil))

}
