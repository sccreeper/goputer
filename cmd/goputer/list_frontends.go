package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sccreeper/goputer/pkg/util"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

type FrontendInfo struct {
	Info struct {
		Name        string `toml:"name"`
		Description string `toml:"description"`
		Authour     string `toml:"authour"`
		Repository  string `toml:"repository"`
		IsPlugin    bool   `toml:"is_plugin"`
	} `toml:"info"`

	Build struct {
		BuildCommand    []string `toml:"command"`
		OutputDirectory string   `toml:"output_dir"`
	} `toml:"build"`
}

func _listFrontends(ctx *cli.Context) error {

	plugin_dir, err := ioutil.ReadDir("./frontends/")
	util.CheckError(err)

	for _, v := range plugin_dir {

		//Load TOML

		tomlFile, err := os.ReadFile(fmt.Sprintf("./frontends/%s/frontend.toml", v.Name()))
		util.CheckError(err)

		var frontendInfo FrontendInfo

		toml.Unmarshal(tomlFile, &frontendInfo)
		util.CheckError(err)

		fmt.Println()
		Bold.Print(frontendInfo.Info.Name + "\n")
		fmt.Println()

		fmt.Printf("%s %s\n", Bold.Sprintf("Description:"), frontendInfo.Info.Description)
		fmt.Printf("%s %s\n", Bold.Sprintf("Authour:"), frontendInfo.Info.Authour)
		fmt.Printf("%s %s\n", Bold.Sprintf("Repository:"), frontendInfo.Info.Repository)

		if frontendInfo.Info.IsPlugin {
			fmt.Printf("%s %s\n", Bold.Sprintf("Is Plugin:"), "Yes")
		} else {
			fmt.Printf("%s %s\n", Bold.Sprintf("Is Plugin:"), "No")
		}

	}

	fmt.Println()

	fmt.Printf("Found %d frontend(s)", len(plugin_dir))

	return nil
}
