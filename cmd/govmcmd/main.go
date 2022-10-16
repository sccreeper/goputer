package main

// VM & Compiler CMD front end

import (
	"errors"
	"log"
	"os"
	"sccreeper/govm/pkg/compiler"
	"sccreeper/govm/pkg/util"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:        "govmcd",
		Description: "Program that is the frontend for running VMs and compiling code",

		Commands: []*cli.Command{
			{
				Name:      "build",
				Aliases:   []string{"b"},
				UsageText: "[file path] [output path]",
				Usage:     "Used to compile programs",
				Action:    _compiler,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func _compiler(ctx *cli.Context) error {

	file_path := ctx.Args().Get(0)

	// See if file exists

	if _, err := os.Stat(file_path); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	log.Printf("Compiling %s", file_path)

	//Read file

	data, err := os.ReadFile(file_path)
	util.CheckError(err)

	err = compiler.Compile(string(data), ctx.Args().Get(1))

	if err != nil {
		util.CheckError(err)
	}

	return nil
}
