//go:build !js && !wasm

package expansions

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"plugin"
	"runtime"
	"sccreeper/goputer/pkg/util"

	"github.com/BurntSushi/toml"
	"golang.org/x/exp/slices"
)

const (
	ExpansionDir          string = "expansions"
	ExpansionManifestName string = "expansion.toml"
)

type ExpansionManifest struct {
	Info struct {
		Name               string   `toml:"name"`
		Description        string   `toml:"description"`
		Authour            string   `toml:"authour"`
		Repository         string   `toml:"repository"`
		Native             bool     `toml:"native"`
		SupportedPlatforms []string `toml:"supported_platforms"`
		ID                 string   `toml:"id"`
		Attributes         []string `toml:"attributes"`
	} `toml:"info"`

	Build struct {
		Command   []string `toml:"command"`
		OutputDir string   `toml:"output_dir"`
	}
}

type ExpansionLoaded struct {
	ExpansionObjectFile *plugin.Plugin
	Handler             func([]byte) []byte
	SetAttribute        func(string, interface{})
	GetAttribute        func(string) interface{}

	Manifest ExpansionManifest
}

var expansions map[string]ExpansionLoaded
var bus_locations map[uint32]string
var native_extension string

func init() {

	if runtime.GOOS == "windows" {
		native_extension = "dll"
	} else {
		native_extension = "so"
	}

	expansions = make(map[string]ExpansionLoaded)
	bus_locations = make(map[uint32]string)

}

// Desktop method for loading expansions
func LoadExpansions() {

	var broken_plugins []string = make([]string, 0)

	directories, err := os.ReadDir(ExpansionDir)
	util.CheckError(err)

	for _, v := range directories {
		log.Printf("Loading %s...\n", v.Name())

		if !v.IsDir() {
			continue
		} else if slices.Contains(broken_plugins, v.Name()) {
			log.Printf("Skipping broken plugin %s\n", v.Name())
			continue
		} else {

			if _, err := os.Stat(path.Join(ExpansionDir, v.Name(), "expansion.toml")); err != nil {
				log.Printf("Plugin manifest for %s is incorrect\n", v.Name())
				log.Println(err.Error())
				continue
			} else {

				var exp_config ExpansionManifest
				var exp_loaded ExpansionLoaded

				// Load TOML
				config_bytes, err := os.ReadFile(path.Join(ExpansionDir, v.Name(), "expansion.toml"))
				loader_error(err, v.Name())

				toml.Unmarshal(config_bytes, &exp_config)
				loader_error(err, v.Name())

				// Load so/lua
				// TODO: Lua
				exp_loaded = ExpansionLoaded{
					Manifest: exp_config,
				}

				if exp_loaded.Manifest.Info.Native {
					// Load symbols from plugin file.

					exp_loaded.ExpansionObjectFile, err = plugin.Open(path.Join(ExpansionDir, v.Name(), fmt.Sprintf("%s.%s", v.Name(), native_extension)))
					loader_error(err, v.Name())

					handler_temp, err := exp_loaded.ExpansionObjectFile.Lookup("Handler")
					if loader_error(err, v.Name()) {
						continue
					}
					exp_loaded.Handler = handler_temp.(func([]byte) []byte)

					get_attribute_temp, err := exp_loaded.ExpansionObjectFile.Lookup("GetAttribute")
					if loader_error(err, v.Name()) {
						continue
					}
					exp_loaded.GetAttribute = get_attribute_temp.(func(string) interface{})

					set_attribute_temp, err := exp_loaded.ExpansionObjectFile.Lookup("SetAttribute")
					if loader_error(err, v.Name()) {
						continue
					}
					exp_loaded.SetAttribute = set_attribute_temp.(func(string, interface{}))

					expansions[exp_loaded.Manifest.Info.ID] = exp_loaded

					log.Printf("Loaded expansion %s successfully.", v.Name())

				} else {
					loader_error(errors.New("lua expansions are not supported yet"), v.Name())
				}

			}

		}

	}

	// Assign all expansions to locations on bus, goputer.sys will be 0

	bus_locations = make(map[uint32]string)
	bus_locations[0] = "goputer.sys"

	var bus_location_index int = 1

	for _, v := range expansions {

		if v.Manifest.Info.ID == "goputer.sys" {
			continue
		} else {
			bus_locations[uint32(bus_location_index)] = v.Manifest.Info.ID
			bus_location_index++
		}

	}

	// Log all bus locations
	log.Println("Bus locations:")

	for k, v := range bus_locations {
		log.Printf("%d: %s\n", k, v)
	}

	//Set bus locations for system module
	for k, v := range bus_locations {
		expansions["goputer.sys"].GetAttribute("expansions").(map[string]uint32)[v] = k
	}

}

func Interaction(location uint32, data []byte) []byte {

	if val, ok := bus_locations[location]; ok {
		return expansions[val].Handler(data)
	} else {
		return []byte{0, 0, 0, 0}
	}

}

// Does a module exist at this location on the bus?
func ModuleExists(location uint32) bool {

	if _, ok := bus_locations[location]; ok {
		return true
	} else {
		return false
	}

}

func SetAttribute(id string, attribute string, value interface{}) {
	expansions[id].SetAttribute(attribute, value)
}

func GetAttribute(id string, attribute string) interface{} {
	return expansions[id].GetAttribute(attribute)
}

func loader_error(e error, expansion_name string) bool {

	if e != nil {
		log.Printf(`
		Error:
		%s
		`, e.Error())
		log.Printf("Failed to load expansion '%s'!\n", expansion_name)

		return true

	} else {
		return false
	}

}
