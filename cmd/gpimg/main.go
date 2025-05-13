package main

import (
	"errors"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"maps"
	"os"
	gpimg "sccreeper/goputer/pkg/gpimg"
	"slices"

	"github.com/urfave/cli/v2"
)

func convertImage(ctx *cli.Context) error {

	// Process flags

	var outputFile *os.File
	var outputPath string = ctx.String("output")
	var inputPath string = ctx.String("input")

	// Open/create files

	_, err := os.Stat(outputPath)
	if err == nil && !ctx.Bool("overwrite") {
		return err
	} else if err == nil && ctx.Bool("overwrite") {
		os.Remove(outputPath)
	}

	outputFile, err = os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Validate and encode flags

	var flags byte
	var flagArgs []string = ctx.StringSlice("format")
	if flagArgs == nil {
		return errors.New("no format flags")
	}

	for _, v := range flagArgs {
		if !slices.Contains(gpimg.FlagNames, v) {
			return fmt.Errorf("unrecognised format flag '%s'", v)
		}
	}

	var hasCompressionFlag bool = false
	for v := range maps.Keys(gpimg.CompressionFlags) {

		var contains bool = slices.Contains(flagArgs, v)
		if contains && hasCompressionFlag {
			return errors.New("can't have more than one compression format")
		} else if contains && !hasCompressionFlag {
			hasCompressionFlag = true
		}
	}

	// Finally encode flags

	for _, v := range flagArgs {

		flags |= gpimg.AllFlags[v]

	}

	gpimg.Encode(inputFile, outputFile, flags)

	return nil
}

func main() {
	app := &cli.App{
		Name:        "gpimg",
		Usage:       "Convert images",
		Description: "Used for converting images to the format used by goputer",
		UsageText:   "gpimg -i in.png -o out.bin --format rle --format opaque",
		Action:      convertImage,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "overwrite",
				Usage: "overwrite destination file if it already exists",
			},
			&cli.StringSliceFlag{
				Name:  "format",
				Usage: "opaque, rle, nocompression",
				Value: cli.NewStringSlice("nocompression", "opaque"),
			},
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Required: true,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
