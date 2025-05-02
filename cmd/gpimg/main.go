package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"sccreeper/goputer/pkg/vm"

	"github.com/urfave/cli/v2"
)

var ErrImgTooLarge error = fmt.Errorf("image is too large must be at most %dx%d", vm.VideoBufferWidth, vm.VideoBufferHeight)

func convertImage(ctx *cli.Context) error {

	// Process flags

	var outputFile *os.File
	var outputPath string = ctx.String("output")
	var inputPath string = ctx.String("input")

	// Open/create files

	_, err := os.Stat(outputPath)
	if errors.Is(err, os.ErrExist) && !ctx.Bool("overwrite") {
		return err
	} else if errors.Is(err, os.ErrExist) && ctx.Bool("overwrite") {
		os.Remove(outputPath)
	}

	outputFile, err = os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Load image

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return err
	}

	if img.Bounds().Dx() > int(vm.VideoBufferWidth) || img.Bounds().Dy() > int(vm.VideoBufferHeight) {
		return ErrImgTooLarge
	}

	// Write width and height as 4 bytes

	var size [4]byte = [4]byte{}

	binary.LittleEndian.PutUint16(
		size[:2],
		uint16(img.Bounds().Dx()),
	)

	binary.LittleEndian.PutUint16(
		size[2:],
		uint16(img.Bounds().Dy()),
	)

	_, err = outputFile.Write(size[:])
	if err != nil {
		return err
	}

	// Write rest of image

	for y := range img.Bounds().Dy() {
		for x := range img.Bounds().Dx() {

			red, green, blue, alpha := img.At(x, y).RGBA()

			_, err = outputFile.Write(
				[]byte{
					byte(red / 257),
					byte(green / 257),
					byte(blue / 257),
					byte(alpha / 257),
				}[:],
			)

			if err != nil {
				return err
			}

		}
	}

	return nil
}

func main() {
	app := &cli.App{
		Name:        "gpimg",
		Usage:       "Convert images",
		Description: "Used for converting images to the format used by goputer",
		UsageText:   "gpimg -i in.png -o out.bin",
		Action:      convertImage,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "overwrite",
				Usage: "overwrite destination file if it already exists",
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
