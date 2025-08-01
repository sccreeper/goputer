package main

import (
	"errors"
	"fmt"
	"os"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/profiler"

	"github.com/urfave/cli/v2"
)

func profile(ctx *cli.Context) error {

	filePath := ctx.String("file")

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	profileData := profiler.Profiler{}

	_, err = profileData.Load(file)
	if err != nil {
		return err
	}

	for _, v := range profileData.Data {

		itn, _ := compiler.DecodeInstructionString(v.Instruction[:])
		fmt.Println(itn)
	}

	return nil

}
