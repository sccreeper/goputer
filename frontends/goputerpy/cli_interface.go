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

	err = os.WriteFile(tmp_file.Name(), program, os.ModePerm)
	util.CheckError(err)

	exec.Command("python3", "py32.py", tmp_file.Name())

}
