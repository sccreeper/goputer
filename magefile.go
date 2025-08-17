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
	"sccreeper/goputer/pkg/expansions"
	"sccreeper/goputer/pkg/util"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type FrontendBuildConfig struct {
	Build struct {
		Command   []string `toml:"command"`
		OutputDir string   `toml:"output_dir"`
		Artifact  string   `toml:"artifact"`
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
func All(includeList string) {

	if _, err := os.Stat("./build"); !os.IsExist(err) {
		os.RemoveAll("./build")
		os.Mkdir("./build", os.ModePerm)
	} else {
		os.Mkdir("./build", os.ModePerm)
	}

	fmt.Println("Building main executable...")

	hash, err := exec.Command("git", "rev-parse", "HEAD").Output()
	util.CheckError(err)

	normalLdFlags := fmt.Sprintf("-X main.Commit=%s", hash)

	previousDir, err := os.Getwd()
	util.CheckError(err)

	os.Chdir("./cmd/goputer/")

	sh.Run("go", "build", "-ldflags", normalLdFlags, "-o", "goputer")

	os.Chdir(previousDir)

	copyFile("./cmd/goputer/goputer", "./build/goputer")

	ldStripFlags := "-s -w"

	//Build launcher

	fmt.Println("Building launcher...")

	previousDir, err = os.Getwd()
	util.CheckError(err)

	os.Chdir("./cmd/gplauncher")

	sh.Run("go", "build", "-ldflags", ldStripFlags, "-o", "gplauncher")

	os.Chdir(previousDir)

	copyFile("./cmd/gplauncher/gplauncher", "./build/gplauncher")

	//Build IDE

	fmt.Println("Building IDE...")

	previousDir, err = os.Getwd()
	util.CheckError(err)

	os.Chdir("./cmd/ide")

	sh.Run("go", "build", "-ldflags", ldStripFlags, "-o", "ide")

	os.Chdir(previousDir)

	copyFile("./cmd/ide/ide", "./build/ide")

	// Build image converter
	fmt.Println("Building gpimg...")

	previousDir, err = os.Getwd()
	util.CheckError(err)

	os.Chdir("./cmd/gpimg")

	sh.Run("go", "build", "-ldflags", ldStripFlags, "-o", "gpimg")

	os.Chdir(previousDir)

	copyFile("./cmd/gpimg/gpimg", "./build/gpimg")

	//Copy the examples

	fmt.Println("Copying examples...")

	os.Mkdir("./build/examples", os.ModePerm)
	os.CopyFS("./build/examples", os.DirFS("./examples/."))

	//Build the frontends
	fmt.Println("Building frontends...")

	os.Mkdir("./build/frontends", os.ModePerm)

	directories, err := ioutil.ReadDir(frontendDir)

	util.CheckError(err)

	for _, v := range directories {

		//Build plugin and copy output folder

		buildToml, err := os.ReadFile(fmt.Sprintf("./frontends/%s/frontend.toml", v.Name()))
		util.CheckError(err)

		var buildConfig FrontendBuildConfig
		toml.Unmarshal(buildToml, &buildConfig)

		if !slices.Contains(strings.Split(includeList, ","), buildConfig.Build.Artifact) {
			fmt.Printf("Skipping artifact %s\n", buildConfig.Build.Artifact)
			continue
		}

		fmt.Printf("Building frontend %s...\n", v.Name())

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
		} else if !slices.Contains(strings.Split(includeList, ","), expConfig.Build.Artifact) {
			fmt.Printf("Skipping artifact %s\n", expConfig.Build.Artifact)
			os.Chdir(previousDir)
			continue
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
func Dev(includeList string) {
	mg.Deps(mg.F(All, includeList))

	compPath, err := os.Getwd()
	util.CheckError(err)

	compPath = filepath.Join(compPath, "build", "goputer")

	previousDir, err := os.Getwd()
	util.CheckError(err)

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

	err = os.Chdir(previousDir)
	util.CheckError(err)

	fmt.Println("Copying testing bins...")

	sh.Copy("./build/profile.gppr", "./bin/profile.gppr")
	sh.Copy("./build/logo.png", "./.github/logo-32.png")

}

func Clean() {
	sh.Rm("./build/")
}

// Build everything, no exclusions
func Build() {
	mg.Deps(mg.F(Dev, "frontend.goputerpy,frontend.gp32,frontend.web,expansion.goputer.sys"))
}
