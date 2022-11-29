package main

import (
	"fmt"
	"os"
	"os/exec"
	"sccreeper/goputer/pkg/util"

	"github.com/BurntSushi/toml"
	"github.com/sqweek/dialog"
)

func open_code() {

	filename, err := dialog.File().Title("Open Program").Load()
	util.CheckError(err)

	CodePath = filename

	LblSelectedCode.SetText(CodePath)

	CodeOpened = true

}

func run_code() {

	cmd := exec.Command("./goputer", "r", "-e", CodePath, "-f", Frontends[SelectedFrontend].PluginSO)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Start()
	util.CheckError(err)

	err = cmd.Process.Release()
	util.CheckError(err)

	App.Quit()

}

func get_frontends() {

	plugin_dir, err := os.ReadDir("./frontends/")
	util.CheckError(err)

	for _, f := range plugin_dir {

		toml_f, err := os.ReadFile(fmt.Sprintf("./frontends/%s/frontend.toml", f.Name()))
		util.CheckError(err)

		var decoded FrontendInfo

		toml.Unmarshal(toml_f, &decoded)

		decoded.PluginSO = f.Name()

		if decoded.Info.IsPlugin {
			Frontends = append(Frontends, decoded)
		}

	}

}
