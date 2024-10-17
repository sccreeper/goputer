//go:build mage

package main

import (
	"fmt"
	"io"
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
	buildDir      string = "./build"
	examplesDir   string = "examples/."
	goputerCmdOut string = "./build/goputer"
	frontendDir   string = "./frontends"
	expansionsDir string = "./expansions/"
)

func copyFile(src string, dest string) error {

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	destFile.Chmod(srcInfo.Mode().Perm())

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil

}

// Builds, clears directory beforehand & copies examples
func All() {

	if _, err := os.Stat("./build"); !os.IsExist(err) {
		os.RemoveAll("./build")
		os.Mkdir("./build", os.ModePerm)
	} else {
		os.Mkdir("./build", os.ModePerm)
	}

	fmt.Println("Building main executable...")

	hash, err := exec.Command("git", "rev-parse", "HEAD").Output()
	util.CheckError(err)

	normal_ldflags := fmt.Sprintf("-s -w -X main.Commit=%s", hash)

	previousDir, err := os.Getwd()
	util.CheckError(err)

	os.Chdir("./cmd/goputer/")

	sh.Run("go", "build", "-ldflags", normal_ldflags, "-o", "goputer")

	os.Chdir(previousDir)

	copyFile("./cmd/goputer/goputer", "./build/goputer")

	//Build launcher

	fmt.Println("Building launcher...")

	launcherLdFlags := "-s -w"

	previousDir, err = os.Getwd()
	util.CheckError(err)

	os.Chdir("./cmd/gplauncher")

	sh.Run("go", "build", "-ldflags", launcherLdFlags, "-o", "gplauncher")

	os.Chdir(previousDir)

	copyFile("./cmd/gplauncher/gplauncher", "./build/gplauncher")

	//Build IDE

	fmt.Println("Building IDE...")

	previousDir, err = os.Getwd()
	util.CheckError(err)

	os.Chdir("./cmd/ide")

	sh.Run("go", "build", "-ldflags", launcherLdFlags, "-o", "ide")

	os.Chdir(previousDir)

	copyFile("./cmd/ide/ide", "./build/ide")

	//Copy the examples

	fmt.Println("Copying examples...")

	sh.Run("mkdir", "./build/examples")
	os.CopyFS("./build/examples", os.DirFS("./examples/."))

	//Build the frontends
	fmt.Println("Building frontends...")

	sh.Run("mkdir", "./build/frontends")

	directories, err := ioutil.ReadDir(frontendDir)

	util.CheckError(err)

	for _, v := range directories {

		fmt.Printf("Building frontend %s...\n", v.Name())

		util.CheckError(err)

		//Build plugin and copy output folder

		buildToml, err := os.ReadFile(fmt.Sprintf("./frontends/%s/frontend.toml", v.Name()))
		util.CheckError(err)

		var buildConfig FrontendBuildConfig
		toml.Unmarshal(buildToml, &buildConfig)

		previousDir, err := os.Getwd()
		util.CheckError(err)

		os.Chdir(fmt.Sprintf("./frontends/%s/", v.Name()))

		cmd := exec.Command(buildConfig.Build.Command[0], buildConfig.Build.Command[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin

		cmd.Run()

		//Escape back to previous directory
		os.Chdir(previousDir)

		os.CopyFS(
			fmt.Sprintf("./build/frontends/%s", v.Name()),
			os.DirFS(fmt.Sprintf("./frontends/%s/%s", v.Name(), buildConfig.Build.OutputDir)),
		)
		copyFile(
			fmt.Sprintf("./frontends/%s/frontend.toml", v.Name()),
			fmt.Sprintf("./build/frontends/%s/frontend.toml", v.Name()),
		)
	}

	// Build expansions

	fmt.Println("Building expansions...")

	directories, err = ioutil.ReadDir(expansionsDir)
	util.CheckError(err)

	for _, v := range directories {

		previousDir, err = os.Getwd()
		util.CheckError(err)

		os.Chdir(path.Join("./expansions", v.Name()))

		var expConfig expansions.ExpansionManifest
		_, err := toml.DecodeFile("./expansion.toml", &expConfig)

		if err != nil {
			fmt.Printf("Config error with expansion %s:\n", v.Name())
			fmt.Println(err.Error())
			continue
		} else if !slices.Contains(expConfig.Info.SupportedPlatforms, runtime.GOOS) {
			fmt.Printf("Error cannot build expansion %s (%s) on current platform %s as it only supports platforms %s\n",
				expConfig.Info.Name,
				expConfig.Info.ID,
				runtime.GOOS,
				strings.Join(expConfig.Info.SupportedPlatforms, ", "))
		}

		fmt.Printf("Building expansion '%s' (%s)...\n", expConfig.Info.Name, expConfig.Info.ID)

		cmd := exec.Command(expConfig.Build.Command[0], expConfig.Build.Command[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin

		cmd.Run()

		os.Chdir(previousDir)

		os.Mkdir("./build/expansions/", os.ModePerm)
		os.Mkdir(fmt.Sprintf("./build/expansions/%s", expConfig.Info.ID), os.ModePerm)

		copyFile(
			fmt.Sprintf("./expansions/%s/expansion.toml", v.Name()),
			fmt.Sprintf("./build/expansions/%s/expansion.toml", expConfig.Info.ID),
		)

		os.CopyFS(
			fmt.Sprintf("./build/expansions/%s/.", expConfig.Info.ID),
			os.DirFS(fmt.Sprintf("./expansions/%s/%s/.", v.Name(), expConfig.Build.OutputDir)),
		)

	}

}

// Builds the examples as well as main All command
func Dev() {
	mg.Deps(All)

	compPath, err := os.Getwd()
	util.CheckError(err)

	compPath = filepath.Join(compPath, "build", "goputer")

	fmt.Println("Building example programs...")

	os.Chdir("./build/examples/")

	dir, err := os.Getwd()
	util.CheckError(err)

	exampleList, err := ioutil.ReadDir(dir)
	util.CheckError(err)

	for _, v := range exampleList {

		if filepath.Ext(v.Name()) != ".gpasm" {
			continue
		} else {
			sh.Run(compPath, "b", v.Name())
		}

	}

}

func Clean() {
	sh.Rm("./build/")
}
