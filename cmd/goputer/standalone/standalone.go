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
var pluginBytes []byte

//go:embed {{.program_code}}
var programCode []byte

var ProgramCheck string
var PluginCheck string

func main() {

	//Decompress
	pluginBytes = decompress(pluginBytes)
	programCode = decompress(programCode)

	//Create temporary plugin and load

	pluginFile, err := os.CreateTemp(os.TempDir(), "*")
	checkErr(err)

	pluginFile.Write(pluginBytes)

	//Verify file integrity

	if ProgramCheck != checksum(programCode) {
		log.Println("Invalid checksum!")
		log.Println(checksum(programCode))
		os.Exit(1)
	}

	frontend_plugin, err := plugin.Open(pluginFile.Name())
	checkErr(err)

	plugin_file_bytes, err := os.ReadFile(pluginFile.Name())
	checkErr(err)

	if PluginCheck != checksum(plugin_file_bytes) || PluginCheck != checksum(pluginBytes) {
		log.Println("Invalid checksum")
		os.Exit(1)
	}

	run_func, err := frontend_plugin.Lookup("Run")
	checkErr(err)

	run_func.(func([]byte, []string))(programCode, []string{filepath.Base(os.Args[0])})

}

func checkErr(err error) {

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
	checkErr(err)
	defer z.Close()

	data, err := io.ReadAll(z)
	util.CheckError(err)

	return data

}
