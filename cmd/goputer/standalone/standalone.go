package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"sccreeper/goputer/pkg/util"
)

//go:embed {{.plugin}}
var plugin_bytes []byte

//go:embed {{.program_code}}
var program_code []byte

var ProgramCheck string
var PluginCheck string

func main() {

	//Decompress
	plugin_bytes = decompress(plugin_bytes)
	program_code = decompress(program_code)

	//Create temporary plugin and load

	plugin_file, err := os.CreateTemp(os.TempDir(), "*")
	check_err(err)

	plugin_file.Write(plugin_bytes)

	//Verify file integrity

	if ProgramCheck != checksum(program_code) {
		log.Println("Invalid checksum!")
		log.Println(checksum(program_code))
		os.Exit(1)
	}

	frontend_plugin, err := plugin.Open(plugin_file.Name())
	check_err(err)

	plugin_file_bytes, err := os.ReadFile(plugin_file.Name())

	if PluginCheck != checksum(plugin_file_bytes) || PluginCheck != checksum(plugin_bytes) {
		log.Println("Invalid checksum")
		os.Exit(1)
	}

	run_func, err := frontend_plugin.Lookup("Run")

	run_func.(func([]byte, []string))(program_code, []string{filepath.Base(os.Args[0])})

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

func decompress(b []byte) []byte {

	z, err := zlib.NewReader(bytes.NewBuffer(b))
	defer z.Close()

	data, err := io.ReadAll(z)
	util.CheckError(err)

	return data

}
