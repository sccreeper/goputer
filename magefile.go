//go:build mage

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sccreeper/govm/pkg/util"

	"github.com/magefile/mage/sh"
)

type FrontendBuildConfig struct {
	Command   []string `json:"command"`
	OutputDir string   `json:"output_dir"`
	IsPlugin  bool     `json:"is_plugin"`
}

var env_map map[string]string = map[string]string{

	"BUILD_DIR":    "./build",
	"FRONTEND_DIR": "./frontends",
}

const (
	build_dir       string = "./build"
	examples_dir    string = "examples/."
	normal_ldflags  string = "-X main.Commit=Help"
	goputer_cmd_out string = "./build/goputer"
	frontend_dir    string = "./frontends"
)

// Builds, clears directory beforehand & copies examples
func All() {

	sh.Rm("./build")
	sh.Run("mkdir", "./build")

	sh.Run("go", "build", "-ldflags", normal_ldflags, "-o", goputer_cmd_out, "./cmd/goputer/main.go")

	//Copy the examples

	sh.Run("mkdir", "./build/examples")
	sh.Run("cp", "-rf", "./examples/.", "./build/examples")

	//Build the frontends

	sh.Run("mkdir", "./build/frontends")

	directories, err := ioutil.ReadDir(frontend_dir)

	util.CheckError(err)

	for _, v := range directories {

		util.CheckError(err)

		//Build plugin and copy output folder

		build_json, err := os.ReadFile(fmt.Sprintf("./frontends/%s/build.json", v.Name()))
		util.CheckError(err)

		var build_config FrontendBuildConfig
		json.Unmarshal(build_json, &build_config)

		previous_dir, err := os.Getwd()
		util.CheckError(err)

		os.Chdir(fmt.Sprintf("./frontends/%s/", v.Name()))

		sh.Run(build_config.Command[0], build_config.Command[1:]...)

		//Escape back to previous directory
		os.Chdir(previous_dir)

		sh.Run("cp", "-rf", fmt.Sprintf("./frontends/%s/build", v.Name()), fmt.Sprintf("./build/frontends/%s", v.Name()))

	}

}
