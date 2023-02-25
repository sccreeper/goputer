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
var native_extension string

func init() {

	if runtime.GOOS == "windows" {
		native_extension = "dll"
	} else {
		native_extension = "so"
	}

}

// Desktop method for loading expansions
func LoadExpansions() {

	var current_plugin string
	var broken_plugins []string = make([]string, 0)

	directories, err := os.ReadDir(ExpansionDir)
	util.CheckError(err)

	defer func() {
		if recover() != nil {
			log.Printf("Fatal error loading the entrypoint file for %s\n.", current_plugin)
			broken_plugins = append(broken_plugins, current_plugin)
		}
	}()

	for _, v := range directories {
		current_plugin = v.Name()

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

					log.Printf("Loaded expansion %s successfully.", v.Name())

				} else {
					loader_error(errors.New("lua expansions are not supported yet"), v.Name())
				}

			}

		}

	}

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
