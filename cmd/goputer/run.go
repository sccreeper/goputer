package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/util"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

// Run the program using a frontend
func runFrontend(ctx *cli.Context) error {

	programBytes, err := os.ReadFile(gpExec)
	util.CheckError(err)

	if string(programBytes[:4]) != compiler.MagicString {
		fmt.Println("Error: Invalid file")
		os.Exit(1)
	}

	// Parse frontend args

	var frontendInfo FrontendInfo

	_, err = toml.DecodeFile(fmt.Sprintf("./frontends/%s/frontend.toml", frontendToUse), &frontendInfo)
	util.CheckError(err)

	if len(frontendInfo.Info.RunCommand) == 0 {
		fmt.Println("Frontend cannot be run as standalone")
		os.Exit(1)
	}

	var formattedCommand []string = make([]string, 0, len(frontendInfo.Info.RunCommand))

	var hasExecutable bool

	for _, v := range frontendInfo.Info.RunCommand {

		// Arguments formatted as $name:value or just $name, e.g. $useprofiler:--useprofiler or just $useprofiler
		// Without a colon gets replaced with value, otherwise whatever is after the colon
		if v[0] == '$' {

			argSplit := strings.Split(v[1:], ":")

			switch argSplit[0] {
			case "use_profiler":
				if useProfiler {
					if len(argSplit) > 1 {
						formattedCommand = append(formattedCommand, argSplit[1])
					} else {
						formattedCommand = append(formattedCommand, "true")
					}
				}
			case "profiler_out":
				if useProfiler {
					if len(argSplit) > 1 {
						formattedCommand = append(formattedCommand, argSplit[1])
					} else {
						formattedCommand = append(formattedCommand, profilerOut)
					}
				}
			case "executable":

				absolutePath, err := filepath.Abs(gpExec)
				if err != nil {
					return err
				}

				formattedCommand = append(formattedCommand, absolutePath)
				hasExecutable = true
			}

		} else {
			formattedCommand = append(formattedCommand, v)
		}
	}

	if !hasExecutable {
		fmt.Println("Frontend must accept executable as argument")
		os.Exit(1)
	}

	cliExecPath, err := os.Executable()
	if err != nil {
		return err
	}

	cliExecPath, err = filepath.EvalSymlinks(cliExecPath)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		formattedCommand[0] += ".exe"
		formattedCommand[0] = strings.ReplaceAll(formattedCommand[0], "/", "\\")
		formattedCommand[0] = filepath.Join(filepath.Dir(cliExecPath), "frontends", frontendToUse, formattedCommand[0])
	}

	cmd := exec.Command(formattedCommand[0], formattedCommand[1:]...)
	cmd.Dir = filepath.Join(filepath.Dir(cliExecPath), "frontends", frontendToUse)

	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("GOPUTER_ROOT=%s", filepath.Dir(cliExecPath)),
		fmt.Sprintf("GOPUTER_CWD=%s", cwd),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil

}
