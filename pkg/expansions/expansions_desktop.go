//go:build !js && !wasm

package expansions

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sccreeper/goputer/pkg/util"
	"slices"

	"github.com/BurntSushi/toml"
	"github.com/Shopify/go-lua"
)

const (
	expansionManifestName string = "expansion.toml"
)

var (
	expansionDir string = "expansions"
)

type ExpansionManifest struct {
	Info struct {
		Name        string `toml:"name"`
		Description string `toml:"description"`
		Authour     string `toml:"authour"`
		Repository  string `toml:"repository"`

		ID         string   `toml:"id"`
		Attributes []string `toml:"attributes"`
	} `toml:"info"`

	Build struct {
		Command   []string `toml:"command"`
		OutputDir string   `toml:"output_dir"`
		Artifact  string   `toml:"artifact"`
	}
}

type ExpansionLoaded struct {
	LuaVM        *lua.State
	Handler      func([]byte) []byte
	SetAttribute func(string, []byte)
	GetAttribute func(string) []byte

	Manifest ExpansionManifest
}

var expansions map[string]ExpansionLoaded
var busLocations map[uint32]string
var nativeExtension string

func init() {

	if runtime.GOOS == "windows" {
		nativeExtension = "dll"
	} else {
		nativeExtension = "so"
	}

	expansions = make(map[string]ExpansionLoaded)
	busLocations = make(map[uint32]string)

}

// Desktop method for loading expansions
func LoadExpansions() {

	var brokenPlugins []string = make([]string, 0)

	expansionDir = filepath.Join(os.Getenv("GOPUTER_ROOT"), expansionDir)

	directories, err := os.ReadDir(expansionDir)
	util.CheckError(err)

	for _, v := range directories {
		log.Printf("Loading %s...\n", v.Name())

		if !v.IsDir() {
			continue
		} else if slices.Contains(brokenPlugins, v.Name()) {
			log.Printf("Skipping broken plugin %s\n", v.Name())
			continue
		} else {

			if _, err := os.Stat(path.Join(expansionDir, v.Name(), expansionManifestName)); err != nil {
				log.Printf("Plugin manifest for %s is incorrect\n", v.Name())
				log.Println(err.Error())
				continue
			} else {

				var expConfig ExpansionManifest
				var expLoaded ExpansionLoaded

				// Load TOML
				configBytes, err := os.ReadFile(path.Join(expansionDir, v.Name(), "expansion.toml"))
				loaderError(err, v.Name())

				toml.Unmarshal(configBytes, &expConfig)
				if loaderError(err, v.Name()) {
					continue
				}

				expLoaded = ExpansionLoaded{
					LuaVM:    lua.NewState(),
					Manifest: expConfig,
				}

				lua.OpenLibraries(expLoaded.LuaVM)
				setStubs(expLoaded.LuaVM)

				err = lua.DoFile(expLoaded.LuaVM, filepath.Join(expansionDir, v.Name(), fmt.Sprintf("%s.lua", expConfig.Info.ID)))
				if loaderError(err, v.Name()) {
					continue
				}

				if !checkForFunction(expLoaded.LuaVM, "Handler") {
					loaderError(errors.New("no handler function"), v.Name())
					continue
				}

				expLoaded.Handler = func(b []byte) []byte {
					expLoaded.LuaVM.Global("Handler")

					expLoaded.LuaVM.NewTable()
					for i, v := range b {
						expLoaded.LuaVM.PushInteger(int(v))
						expLoaded.LuaVM.RawSetInt(-2, i+1)
					}

					expLoaded.LuaVM.Call(1, 1)

					return getBytesFromStack(expLoaded.LuaVM)

				}

				if !checkForFunction(expLoaded.LuaVM, "SetAttribute") {
					loaderError(errors.New("no set attribute function"), v.Name())
					continue
				}

				expLoaded.SetAttribute = func(s string, data []byte) {
					expLoaded.LuaVM.Global("SetAttribute")
					expLoaded.LuaVM.PushString(s)
					expLoaded.LuaVM.NewTable()

					for i, v := range data {
						expLoaded.LuaVM.PushInteger(int(v))
						expLoaded.LuaVM.RawSetInt(-2, i+1)
					}

					expLoaded.LuaVM.Call(2, 0)
				}

				if !checkForFunction(expLoaded.LuaVM, "GetAttribute") {
					loaderError(errors.New("no get attribute function"), v.Name())
					continue
				}

				expLoaded.GetAttribute = func(s string) []byte {
					expLoaded.LuaVM.Global("GetAttribute")
					expLoaded.LuaVM.PushString(s)
					expLoaded.LuaVM.Call(1, 1)

					return getBytesFromStack(expLoaded.LuaVM)

				}

				expansions[v.Name()] = expLoaded

			}

		}

	}

	// Assign all expansions to locations on bus, goputer.sys will be 0

	busLocations = make(map[uint32]string)
	busLocations[0] = "goputer.sys"

	var busLocationIndex int = 1

	for _, v := range expansions {

		if v.Manifest.Info.ID == "goputer.sys" {
			continue
		} else {
			busLocations[uint32(busLocationIndex)] = v.Manifest.Info.ID
			busLocationIndex++
		}

	}

	// Log all bus locations
	log.Println("Bus locations:")

	for k, v := range busLocations {
		log.Printf("%d: %s\n", k, v)
	}

	busLocationBytes := make([]byte, 0)

	//Set bus locations for system module
	for k, v := range busLocations {
		busLocationBytes = append(busLocationBytes, []byte(fmt.Sprintf("%s\x00%d\x00", v, k))...)
	}

}

func Interaction(location uint32, data []byte) []byte {

	if val, ok := busLocations[location]; ok {
		return expansions[val].Handler(data)
	} else {
		return []byte{0, 0, 0, 0}
	}

}

// Does a module exist at this location on the bus?
func ModuleExists(location uint32) bool {

	if _, ok := busLocations[location]; ok {
		return true
	} else {
		return false
	}

}

func SetAttribute(id string, attribute string, value []byte) {
	expansions[id].SetAttribute(attribute, value)
}

func GetAttribute(id string, attribute string) []byte {
	return expansions[id].GetAttribute(attribute)
}

func loaderError(e error, expansionName string) bool {

	if e != nil {
		log.Printf("Error: %s\n", e.Error())
		log.Printf("Failed to load expansion '%s'!\n", expansionName)

		return true

	} else {
		return false
	}

}

func checkForFunction(vm *lua.State, name string) bool {
	vm.Global(name)
	if !vm.IsFunction(-1) {
		return false
	}
	vm.Pop(1)

	return true
}

func getBytesFromStack(vm *lua.State) (result []byte) {
	result = make([]byte, 0)
	vm.Length(-1)
	length, _ := vm.ToInteger(-1)
	vm.Pop(1)

	for i := 1; i <= length; i++ {
		vm.RawGetInt(-1, i)
		if num, ok := vm.ToInteger(-1); ok {
			result = append(result, byte(num))
		}
		vm.Pop(1)
	}

	vm.Pop(1)
	return result
}

func setStubs(vm *lua.State) {

	vm.NewTable()

	vm.PushGoFunction(func(state *lua.State) int {

		if !vm.IsNumber(1) {
			vm.PushNil()
			return 1
		}

		num, _ := vm.ToInteger(1)

		var data [4]byte = [4]byte{}
		binary.LittleEndian.PutUint32(data[:], uint32(num))

		vm.NewTable()
		for i, v := range data {
			vm.PushInteger(int(v))
			vm.RawSetInt(-2, i+1)
		}

		return 1

	})
	vm.SetField(-2, "toLittleEndian")

	vm.PushGoFunction(func(state *lua.State) int {

		if !vm.IsTable(1) {
			vm.PushNil()
			return 1
		}

		numBytes := getBytesFromStack(vm)

		if len(numBytes) != 4 {
			vm.PushNil()
			return 1
		}

		result := binary.LittleEndian.Uint32(numBytes)

		vm.PushInteger(int(result))

		return 1

	})
	vm.SetField(-2, "fromLittleEndian")

	vm.SetGlobal("Gp")

}
