package main // import "github.com/sccreeper/goputer"

// VM & Compiler CMD front end

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	_ "embed"
)

func format_instruction(i_name string, i_data []string) string {

	return fmt.Sprintf("%s %s", color.GreenString(i_name), color.CyanString(strings.Join(i_data[:], " ")))

}

func main() {

	if runtime.GOOS == "windows" {
		PluginExt = ".dll"
	} else {
		PluginExt = ".so"
	}

	app := &cli.App{
		Name:        "goputer",
		Usage:       "CLI for goputer",
		Description: "CLI tool for goputer. Main purpose is running and compiling code.",

		Commands: []*cli.Command{
			{
				Name:    "build",
				Aliases: []string{"b"},
				Usage:   "Used to compile programs",
				Action:  _compiler,

				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "json",
						Usage:       "Enable JSON outputting",
						Destination: &UseJson,
					},
					&cli.StringFlag{
						Name:        "jsonpath",
						Usage:       "Output program structure/data in `FILE` ",
						Destination: &JsonPath,
					},
					&cli.StringFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						Usage:       "Output binary to `FILE`",
						Destination: &OutputPath,
					},
					&cli.BoolFlag{
						Name:        "verbose",
						Aliases:     []string{"v"},
						Usage:       "Verbose log output",
						Destination: &Verbose,
					},
					&cli.BoolFlag{
						Name:        "standalone",
						Aliases:     []string{"s"},
						Usage:       "Create a standalone executable",
						Destination: &IsStandalone,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "frontend",
						Aliases:     []string{"f"},
						Usage:       "Frontend to create standalone with",
						Destination: &FrontendToUse,
						Required:    false,
					},
				},
			},
			{
				Name:    "disassemble",
				Aliases: []string{"d"},
				Usage:   "Used to disassemble programs",
				Action:  _disassemble,
			},
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "Run programs",
				Action:  _run,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "frontend",
						Aliases: []string{"f"},
						Usage:   "Frontend to use",
						// DefaultText: "gp32",
						Destination: &FrontendToUse,
					},
					&cli.StringFlag{
						Name:        "exec",
						Aliases:     []string{"e"},
						Usage:       "Executable to run",
						Destination: &GPExec,
					},
				},
			},
			{
				Name:   "list",
				Usage:  "Lists plugins available",
				Action: _list_frontends,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
