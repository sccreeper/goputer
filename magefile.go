//go:build mage

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sccreeper/goputer/pkg/util"

	"github.com/magefile/mage/mg"
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
	goputer_cmd_out string = "./build/goputer"
	frontend_dir    string = "./frontends"
)

// Builds, clears directory beforehand & copies examples
func All() {

	if _, err := os.Stat("./build"); !os.IsExist(err) {
		sh.Rm("./build")
		sh.Run("mkdir", "./build")
	} else {
		sh.Run("mkdir", "./build")
	}

	fmt.Println("Building main executable...")

	hash, err := exec.Command("git", "rev-parse", "HEAD").Output()
	util.CheckError(err)

	normal_ldflags := fmt.Sprintf("-s -w -X main.Commit=%s", hash)

	sh.Run("go", "build", "-ldflags", normal_ldflags, "-o", goputer_cmd_out, "./cmd/goputer/main.go")

	//Copy the examples

	fmt.Println("Copying examples...")

	sh.Run("mkdir", "./build/examples")
	sh.Run("cp", "-rf", "./examples/.", "./build/examples")

	//Build the frontends
	fmt.Println("Building frontends...")

	sh.Run("mkdir", "./build/frontends")

	directories, err := ioutil.ReadDir(frontend_dir)

	util.CheckError(err)

	for _, v := range directories {

		fmt.Printf("Building frontend %s...\n", v.Name())

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

// Builds the examples as well as main All command
func Dev() {
	mg.Deps(All)

	fmt.Println("Building example programs...")

	os.Chdir("./build/")

	example_list, err := ioutil.ReadDir("./examples")
	util.CheckError(err)

	for _, v := range example_list {

		sh.Run("./goputer", "b", fmt.Sprintf("./examples/%s", v.Name()))

	}

}

func Clean() {
	sh.Rm("./build/")
}
