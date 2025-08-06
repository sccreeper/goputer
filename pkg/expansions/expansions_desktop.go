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
	"slices"

	"github.com/BurntSushi/toml"
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
		Artifact  string   `toml:"artifact"`
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

	directories, err := os.ReadDir(ExpansionDir)
	util.CheckError(err)

	for _, v := range directories {
		log.Printf("Loading %s...\n", v.Name())

		if !v.IsDir() {
			continue
		} else if slices.Contains(brokenPlugins, v.Name()) {
			log.Printf("Skipping broken plugin %s\n", v.Name())
			continue
		} else {

			if _, err := os.Stat(path.Join(ExpansionDir, v.Name(), "expansion.toml")); err != nil {
				log.Printf("Plugin manifest for %s is incorrect\n", v.Name())
				log.Println(err.Error())
				continue
			} else {

				var expConfig ExpansionManifest
				var expLoaded ExpansionLoaded

				// Load TOML
				configBytes, err := os.ReadFile(path.Join(ExpansionDir, v.Name(), "expansion.toml"))
				loaderError(err, v.Name())

				toml.Unmarshal(configBytes, &expConfig)
				loaderError(err, v.Name())

				// Load so/lua
				// TODO: Lua
				expLoaded = ExpansionLoaded{
					Manifest: expConfig,
				}

				if expLoaded.Manifest.Info.Native {
					// Load symbols from plugin file.

					expLoaded.ExpansionObjectFile, err = plugin.Open(path.Join(ExpansionDir, v.Name(), fmt.Sprintf("%s.%s", v.Name(), nativeExtension)))
					loaderError(err, v.Name())

					handlerTemp, err := expLoaded.ExpansionObjectFile.Lookup("Handler")
					if loaderError(err, v.Name()) {
						continue
					}
					expLoaded.Handler = handlerTemp.(func([]byte) []byte)

					getAttributeTemp, err := expLoaded.ExpansionObjectFile.Lookup("GetAttribute")
					if loaderError(err, v.Name()) {
						continue
					}
					expLoaded.GetAttribute = getAttributeTemp.(func(string) interface{})

					setAttributeTemp, err := expLoaded.ExpansionObjectFile.Lookup("SetAttribute")
					if loaderError(err, v.Name()) {
						continue
					}
					expLoaded.SetAttribute = setAttributeTemp.(func(string, interface{}))

					expansions[expLoaded.Manifest.Info.ID] = expLoaded

					log.Printf("Loaded expansion %s successfully.", v.Name())

				} else {
					loaderError(errors.New("lua expansions are not supported yet"), v.Name())
				}

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

	//Set bus locations for system module
	for k, v := range busLocations {
		expansions["goputer.sys"].GetAttribute("expansions").(map[string]uint32)[v] = k
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

func SetAttribute(id string, attribute string, value interface{}) {
	expansions[id].SetAttribute(attribute, value)
}

func GetAttribute(id string, attribute string) interface{} {
	return expansions[id].GetAttribute(attribute)
}

func loaderError(e error, expansionName string) bool {

	if e != nil {
		log.Printf(`
		Error:
		%s
		`, e.Error())
		log.Printf("Failed to load expansion '%s'!\n", expansionName)

		return true

	} else {
		return false
	}

}
