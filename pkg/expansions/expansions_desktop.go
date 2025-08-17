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
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"
	"sccreeper/goputer/pkg/vm"
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
	registryReferences = make(map[string]map[string]int)

	for _, v := range vm.HookNames {
		registryReferences[v] = make(map[string]int)
	}

}

// Desktop method for loading expansions
func LoadExpansions(vm *vm.VM) {

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
				setStubs(expLoaded.LuaVM, vm)

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

					return getBytesFromStack(expLoaded.LuaVM, -1, true)

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

					return getBytesFromStack(expLoaded.LuaVM, -1, true)

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

func checkForFunction(l *lua.State, name string) bool {
	l.Global(name)
	if !l.IsFunction(-1) {
		return false
	}
	l.Pop(1)

	return true
}

func getBytesFromStack(l *lua.State, index int, pop bool) (result []byte) {
	result = make([]byte, 0)
	l.Length(index)
	length, _ := l.ToInteger(-1)
	l.Pop(1)

	for i := 1; i <= length; i++ {
		l.RawGetInt(index, i)
		if num, ok := l.ToInteger(-1); ok {
			result = append(result, byte(num))
		}
		l.Pop(1)
	}

	if pop {
		l.Pop(1)
	}
	return result
}

var nextListenerRef int
var registryReferences map[string]map[string]int

func setStubs(l *lua.State, gpVm *vm.VM) {

	l.NewTable()

	l.PushGoFunction(func(state *lua.State) int {

		if !l.IsNumber(1) {
			l.PushNil()
			return 1
		}

		num, _ := l.ToInteger(1)
		l.Pop(1)

		var data [4]byte = [4]byte{}
		binary.LittleEndian.PutUint32(data[:], uint32(num))

		l.NewTable()
		for i, v := range data {
			l.PushInteger(int(v))
			l.RawSetInt(-2, i+1)
		}

		return 1

	})
	l.SetField(-2, "toLittleEndian")

	l.PushGoFunction(func(state *lua.State) int {

		if !l.IsTable(1) {
			l.PushNil()
			return 1
		}

		numBytes := getBytesFromStack(l, -1, true)

		if len(numBytes) != 4 {
			l.PushNil()
			return 1
		}

		result := binary.LittleEndian.Uint32(numBytes)

		l.PushInteger(int(result))

		return 1

	})
	l.SetField(-2, "fromLittleEndian")

	l.PushGoFunction(func(state *lua.State) int {

		if !l.IsNumber(1) || !l.IsNumber(2) {
			l.PushBoolean(false)
			return 1
		}

		addr, _ := l.ToInteger(1)
		val, _ := l.ToInteger(2)
		l.Pop(2)

		if addr < 0 || addr >= len(gpVm.MemArray) {
			l.PushBoolean(false)
			return 1
		}

		gpVm.MemArray[addr] = byte(val)

		l.PushBoolean(true)
		return 1

	})
	l.SetField(-2, "setMemoryAddress")

	l.PushGoFunction(func(state *lua.State) int {
		if !l.IsNumber(1) {
			l.PushNil()
			return 1
		}

		addr, _ := l.ToInteger(1)
		l.Pop(1)

		if addr < 0 || addr >= len(gpVm.MemArray) {
			l.PushNil()
			return 1
		}

		l.PushInteger(int(gpVm.MemArray[addr]))
		return 1

	})
	l.SetField(-2, "getMemoryAddress")

	l.PushGoFunction(func(state *lua.State) int {
		if !l.IsNumber(1) || !l.IsTable(2) {
			l.PushBoolean(false)
			return 1
		}

		addr, _ := l.ToInteger(1)
		data := getBytesFromStack(l, -1, true)

		if addr < 0 || addr+len(data) > len(gpVm.MemArray) {
			l.PushBoolean(false)
			return 1
		}

		copy(
			gpVm.MemArray[addr:addr+len(data)],
			data[:],
		)

		l.PushBoolean(true)
		return 1

	})
	l.SetField(-2, "setMemoryRange")

	l.PushGoFunction(func(state *lua.State) int {
		if !l.IsNumber(1) || !l.IsNumber(2) {
			l.PushNil()
			return 1
		}

		addr, _ := l.ToInteger(1)
		size, _ := l.ToInteger(2)
		l.Pop(2)

		if addr < 0 || addr+size > len(gpVm.MemArray) || size <= 0 {
			l.PushNil()
			return 1
		}

		l.NewTable()
		for i := range size {
			l.PushInteger(int(gpVm.MemArray[addr+i]))
			l.RawSetInt(-2, i+1)
		}

		return 1
	})
	l.SetField(-2, "getMemoryRange")

	l.PushGoFunction(func(state *lua.State) int {
		if !l.IsNumber(1) || !l.IsNumber(2) || !l.IsNumber(3) {
			l.PushBoolean(false)
			return 1
		}

		addr, _ := l.ToInteger(1)
		size, _ := l.ToInteger(2)
		value, _ := l.ToInteger(3)
		l.Pop(3)

		if addr < 0 || addr+size > len(gpVm.MemArray) {
			l.PushBoolean(false)
			return 1
		}

		for i := addr; i < addr+size; i++ {
			gpVm.MemArray[i] = byte(value)
		}

		l.PushBoolean(true)
		return 1

	})
	l.SetField(-2, "clearMemoryRange")

	l.PushGoFunction(func(state *lua.State) int {

		if !l.IsString(1) {
			l.PushNil()
			return 1
		}

		reg, _ := l.ToString(1)
		l.Pop(1)

		if _, exists := constants.RegisterInts[reg]; !exists {
			l.PushNil()
			return 1
		}

		l.PushInteger(int(gpVm.Registers[constants.RegisterInts[reg]]))
		return 1

	})
	l.SetField(-2, "getRegister")

	l.PushGoFunction(func(state *lua.State) int {
		if !l.IsString(1) || !l.IsNumber(2) {
			l.PushBoolean(false)
			return 1
		}

		reg, _ := l.ToString(1)
		val, _ := l.ToInteger(2)

		if _, exists := constants.RegisterInts[reg]; !exists {
			l.PushBoolean(false)
			return 1
		}

		gpVm.Registers[constants.RegisterInts[reg]] = uint32(val)

		l.PushBoolean(true)
		return 1

	})
	l.SetField(-2, "setRegister")

	l.PushGoFunction(func(state *lua.State) int {
		if !l.IsString(1) || !l.IsTable(2) || !l.IsNumber(3) {
			l.PushBoolean(false)
			return 1
		}

		buf, _ := l.ToString(1)
		val := getBytesFromStack(l, 2, false)
		offset, _ := l.ToInteger(3)
		l.Pop(3)

		if buf != "d0" && buf != "vt" {
			l.PushBoolean(false)
			return 1
		} else if offset+len(val) > 128 {
			l.PushBoolean(false)
			return 1
		}

		if buf == "d0" {
			for i := range val {
				gpVm.DataBuffer[offset+i] = val[i]
			}
		} else {
			for i := range val {
				gpVm.TextBuffer[offset+i] = val[i]
			}
		}

		l.PushBoolean(true)
		return 1

	})
	l.SetField(-2, "setBuffer")

	l.PushGoFunction(func(state *lua.State) int {
		if !l.IsString(1) || !l.IsNumber(2) || !l.IsNumber(3) {
			l.PushNil()
			return 1
		}

		buf, _ := l.ToString(1)
		offset, _ := l.ToInteger(2)
		length, _ := l.ToInteger(3)
		l.Pop(3)

		if buf != "d0" && buf != "vt" {
			l.PushNil()
			return 1
		} else if offset+length > 128 {
			l.PushNil()
			return 1
		}

		var data []byte
		if buf == "d0" {
			data = gpVm.DataBuffer[offset : offset+length]
		} else {
			data = gpVm.TextBuffer[offset : offset+length]
		}

		l.NewTable()
		for i, v := range data {
			l.PushInteger(int(v))
			l.RawSetInt(-2, i+1)
		}

		return 1
	})
	l.SetField(-2, "getBuffer")

	l.PushGoFunction(func(state *lua.State) int {
		gpVm.Finished = true
		return 0
	})
	l.SetField(-2, "stop")

	l.PushGoFunction(func(state *lua.State) int {
		if !l.IsNumber(1) {
			l.PushNil()
			return 1
		}

		addr, _ := l.ToInteger(1)
		l.Pop(1)

		if addr < 0 || addr+int(compiler.InstructionLength) > len(gpVm.MemArray) {
			l.PushNil()
			return 1
		}

		res, err := compiler.DecodeInstructionString(gpVm.MemArray[addr : addr+int(compiler.InstructionLength)])
		if err != nil {
			l.PushNil()
			return 1
		}

		l.PushString(res)
		return 1

	})
	l.SetField(-2, "decodeInstructionToString")

	l.PushGoFunction(func(state *lua.State) int {
		val, ok := l.ToString(1)

		if !ok {
			return 0
		} else {
			log.Println(val)
		}

		return 0
	})
	l.SetField(-2, "log")

	// Hooks

	l.NewTable()

	l.PushGoFunction(func(state *lua.State) int {

		if !l.IsString(1) || !l.IsString(2) || !l.IsFunction(3) {
			log.Println("incorrect args when adding hook")
			l.PushBoolean(false)
			return 1
		}

		name, _ := l.ToString(1)
		event, _ := l.ToString(2)

		ref := nextListenerRef

		l.PushValue(3)
		l.RawSetInt(lua.RegistryIndex, ref)

		l.Pop(3)

		if !slices.Contains(vm.HookNames, event) {
			log.Println("unable to add hook")
			l.PushBoolean(false)
			return 1
		}

		err := gpVm.AddHook(
			name,
			vm.VMHook(slices.Index(vm.HookNames, event)),
			func() {

				l.RawGetInt(lua.RegistryIndex, ref)
				l.Call(0, 0)

			},
		)

		if err != nil {
			log.Printf("unable to add hook: %s\n", err.Error())
			l.PushBoolean(false)
			return 1
		}

		registryReferences[event][name] = nextListenerRef

		nextListenerRef++

		l.PushBoolean(true)
		return 1

	})
	l.SetField(-2, "addHook")

	l.PushGoFunction(func(state *lua.State) int {

		if !l.IsString(1) || !l.IsString(2) {
			l.PushBoolean(false)
			return 1
		}

		name, _ := l.ToString(1)
		event, _ := l.ToString(2)
		l.Pop(2)

		if !slices.Contains(vm.HookNames, event) {
			l.PushBoolean(false)
			return 1
		}

		gpVm.RemoveHook(name, vm.VMHook(slices.Index(vm.HookNames, event)))

		l.PushNil()
		l.RawSetInt(lua.RegistryIndex, registryReferences[event][name])
		delete(registryReferences[event], name)

		l.PushBoolean(true)
		return 1

	})
	l.SetField(-2, "removeHook")

	l.SetField(-2, "hooks")

	l.SetGlobal("Gp")

}
