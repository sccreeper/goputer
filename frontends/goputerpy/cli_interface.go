// Plugin glue between goputer CLI and Python
package main

import (
	"os"
	"os/exec"
	"sccreeper/goputer/pkg/util"
)

func Run(program []byte, args []string) {

	//Write program bytes to tmp file and then call command to run Python.

	tmp_file, err := os.CreateTemp(os.TempDir(), "")
	util.CheckError(err)

	err = os.WriteFile(tmp_file.Name(), program, 0555)
	util.CheckError(err)

	cmd := exec.Command("./frontends/goputerpy/goputerpy", tmp_file.Name())
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Start()

}
