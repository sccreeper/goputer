package main // import "github.com/sccreeper/goputer"

// VM & Compiler CMD front end

import (
	"log"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"
)

func main() {

	if runtime.GOOS == "windows" {
		pluginExt = ".dll"
	} else {
		pluginExt = ".so"
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
						Destination: &useJson,
					},
					&cli.StringFlag{
						Name:        "jsonpath",
						Usage:       "Output program structure/data in `FILE` ",
						Destination: &jsonPath,
					},
					&cli.StringFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						Usage:       "Output binary to `FILE`",
						Destination: &programCompileOut,
					},
					&cli.BoolFlag{
						Name:        "verbose",
						Aliases:     []string{"v"},
						Usage:       "Verbose log output",
						Destination: &beVerbose,
					},
					&cli.BoolFlag{
						Name:        "standalone",
						Aliases:     []string{"s"},
						Usage:       "Create a standalone executable",
						Destination: &isStandalone,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "frontend",
						Aliases:     []string{"f"},
						Usage:       "Frontend to create standalone with",
						Destination: &frontendToUse,
						Required:    false,
					},
				},
			},
			{
				Name:    "disassemble",
				Aliases: []string{"d"},
				Usage:   "Used to disassemble programs",
				Action:  _disassemble,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "verbose",
						Aliases:     []string{"v"},
						Destination: &beVerbose,
					},
				},
			},
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "Run programs",
				Action:  runFrontend,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "frontend",
						Aliases: []string{"f"},
						Usage:   "Frontend to use",
						// DefaultText: "gp32",
						Destination: &frontendToUse,
					},
					&cli.StringFlag{
						Name:        "exec",
						Aliases:     []string{"e"},
						Usage:       "Executable to run",
						Destination: &gpExec,
					},
					&cli.BoolFlag{
						Name:        "useprofiler",
						Usage:       "Wether or not to use the profiler",
						Destination: &useProfiler,
					},
					&cli.StringFlag{
						Name:        "profilerout",
						Usage:       "Output profiler to `FILE`",
						Destination: &profilerOut,
					},
				},
			},
			{
				Name:    "profile",
				Aliases: []string{"p"},
				Usage:   "Load a profile (.gppr) file",
				Action:  profile,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "File to load",
						Required: true,
					},
				},
			},
			{
				Name:   "list",
				Usage:  "Lists plugins available",
				Action: listFrontends,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
