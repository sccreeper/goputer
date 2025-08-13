package main

import (
	"fmt"
	"os"
	"sccreeper/goputer/pkg/util"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

type FrontendInfo struct {
	Info struct {
		Name        string   `toml:"name"`
		Description string   `toml:"description"`
		Authour     string   `toml:"authour"`
		Repository  string   `toml:"repository"`
		RunCommand  []string `toml:"run"`
	} `toml:"info"`

	Build struct {
		BuildCommand    []string `toml:"command"`
		OutputDirectory string   `toml:"output_dir"`
	} `toml:"build"`
}

func listFrontends(ctx *cli.Context) error {

	pluginDir, err := os.ReadDir("./frontends/")
	util.CheckError(err)

	for _, v := range pluginDir {

		//Load TOML

		tomlFile, err := os.ReadFile(fmt.Sprintf("./frontends/%s/frontend.toml", v.Name()))
		util.CheckError(err)

		var frontendInfo FrontendInfo

		toml.Unmarshal(tomlFile, &frontendInfo)
		util.CheckError(err)

		fmt.Println()
		bold.Print(frontendInfo.Info.Name + "\n")
		fmt.Println()

		fmt.Printf("%s %s\n", bold.Sprintf("Description:"), frontendInfo.Info.Description)
		fmt.Printf("%s %s\n", bold.Sprintf("Authour:"), frontendInfo.Info.Authour)
		fmt.Printf("%s %s\n", bold.Sprintf("Repository:"), frontendInfo.Info.Repository)

		if len(frontendInfo.Info.RunCommand) > 0 {
			fmt.Printf("%s %s\n", bold.Sprintf("Is executable?:"), "Yes")
		} else {
			fmt.Printf("%s %s\n", bold.Sprintf("Is executable?:"), "No")
		}

	}

	fmt.Println()

	fmt.Printf("Found %d frontend(s)", len(pluginDir))

	return nil
}
