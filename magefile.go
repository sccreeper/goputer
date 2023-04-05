//go:build mage

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sccreeper/goputer/pkg/expansions"
	"sccreeper/goputer/pkg/util"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"golang.org/x/exp/slices"
)

type FrontendBuildConfig struct {
	Build struct {
		Command   []string `toml:"command"`
		OutputDir string   `toml:"output_dir"`
	} `toml:"build"`
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
	expansions_dir  string = "./expansions/"
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

	previous_dir, err := os.Getwd()
	util.CheckError(err)

	os.Chdir("./cmd/goputer/")

	sh.Run("go", "build", "-ldflags", normal_ldflags, "-o", "goputer")

	os.Chdir(previous_dir)

	sh.Run("cp", "./cmd/goputer/goputer", "./build/goputer")

	//Build launcher

	fmt.Println("Building launcher...")

	launcher_ld_flags := "-s -w"

	previous_dir, err = os.Getwd()
	util.CheckError(err)

	os.Chdir("./cmd/gplauncher")

	sh.Run("go", "build", "-ldflags", launcher_ld_flags, "-o", "gplauncher")

	os.Chdir(previous_dir)

	sh.Run("cp", "./cmd/gplauncher/gplauncher", "./build/gplauncher")

	//Build IDE

	fmt.Println("Building IDE...")

	previous_dir, err = os.Getwd()
	util.CheckError(err)

	os.Chdir("./cmd/ide")

	sh.Run("go", "build", "-ldflags", launcher_ld_flags, "-o", "ide")

	os.Chdir(previous_dir)

	sh.Run("cp", "./cmd/ide/ide", "./build/ide")

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

		build_toml, err := os.ReadFile(fmt.Sprintf("./frontends/%s/frontend.toml", v.Name()))
		util.CheckError(err)

		var build_config FrontendBuildConfig
		toml.Unmarshal(build_toml, &build_config)

		previous_dir, err := os.Getwd()
		util.CheckError(err)

		os.Chdir(fmt.Sprintf("./frontends/%s/", v.Name()))

		cmd := exec.Command(build_config.Build.Command[0], build_config.Build.Command[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin

		cmd.Run()

		//Escape back to previous directory
		os.Chdir(previous_dir)

		sh.Run("cp", "-rf", fmt.Sprintf("./frontends/%s/%s", v.Name(), build_config.Build.OutputDir), fmt.Sprintf("./build/frontends/%s", v.Name()))
		sh.Run("cp", fmt.Sprintf("./frontends/%s/frontend.toml", v.Name()), fmt.Sprintf("./build/frontends/%s/frontend.toml", v.Name()))
	}

	// Build expansions

	fmt.Println("Building expansions...")

	directories, err = ioutil.ReadDir(expansions_dir)
	util.CheckError(err)

	for _, v := range directories {

		previous_dir, err = os.Getwd()
		util.CheckError(err)

		os.Chdir(path.Join("./expansions", v.Name()))

		var exp_config expansions.ExpansionManifest
		_, err := toml.DecodeFile("./expansion.toml", &exp_config)

		if err != nil {
			fmt.Printf("Config error with expansion %s:\n", v.Name())
			fmt.Println(err.Error())
			continue
		} else if !slices.Contains(exp_config.Info.SupportedPlatforms, runtime.GOOS) {
			fmt.Printf("Error cannot build expansion %s (%s) on current platform %s as it only supports platforms %s\n",
				exp_config.Info.Name,
				exp_config.Info.ID,
				runtime.GOOS,
				strings.Join(exp_config.Info.SupportedPlatforms, ", "))
		}

		fmt.Printf("Building expansion '%s' (%s)...\n", exp_config.Info.Name, exp_config.Info.ID)

		cmd := exec.Command(exp_config.Build.Command[0], exp_config.Build.Command[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin

		cmd.Run()

		os.Chdir(previous_dir)

		os.Mkdir("./build/expansions/", 0755)
		os.Mkdir(fmt.Sprintf("./build/expansions/%s", exp_config.Info.ID), 0755)

		sh.Run("cp", fmt.Sprintf("./expansions/%s/expansion.toml", v.Name()), fmt.Sprintf("./build/expansions/%s/expansion.toml", exp_config.Info.ID))
		sh.Run("cp", "-ra", fmt.Sprintf("./expansions/%s/%s/.", v.Name(), exp_config.Build.OutputDir), fmt.Sprintf("./build/expansions/%s/", exp_config.Info.ID))

	}

}

// Builds the examples as well as main All command
func Dev() {
	mg.Deps(All)

	comp_path, err := os.Getwd()
	util.CheckError(err)

	comp_path += "/build/goputer"

	fmt.Println("Building example programs...")

	os.Chdir("./build/examples/")

	dir, err := os.Getwd()
	util.CheckError(err)

	example_list, err := ioutil.ReadDir(dir)
	util.CheckError(err)

	for _, v := range example_list {

		if filepath.Ext(v.Name()) != ".gpasm" {
			continue
		}

		sh.Run(comp_path, "b", v.Name())

	}

}

func Clean() {
	sh.Rm("./build/")
}
