package main

import (
	"errors"
	"fmt"
	"os"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/profiler"
	"slices"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

	profileEntriesSorted := make([]profiler.ProfileEntry, 0, len(profileData.Data))

	for _, v := range profileData.Data {
		profileEntriesSorted = append(profileEntriesSorted, *v)
	}

	slices.SortFunc(profileEntriesSorted, func(a profiler.ProfileEntry, b profiler.ProfileEntry) int {
		if a.Address < b.Address {
			return -1
		} else if a.Address > b.Address {
			return 1
		} else {
			return 0
		}
	})

	// Setup terminal app

	app := tview.NewApplication().SetTitle("goputer Profiler")

	table := tview.NewTable()
	table.SetTitle("Instructions executed")

	for r, v := range profileEntriesSorted {

		table.SetCell(
			r, 0,
			tview.NewTableCell(
				fmt.Sprintf("0x%08X", v.Address),
			).SetTextColor(
				tcell.ColorGreen,
			).SetAlign(tview.AlignCenter),
		)

		itnString, err := compiler.DecodeInstructionString(v.Instruction[:])
		if err != nil {
			panic(err)
		}

		table.SetCell(
			r, 1,
			tview.NewTableCell(
				itnString,
			).SetMaxWidth(24),
		)

		table.SetCell(
			r, 2,
			tview.NewTableCell(
				strconv.FormatInt(int64(v.TotalTimesExecuted), 10),
			).SetAlign(tview.AlignRight),
		)

		table.SetCell(
			r, 3,
			tview.NewTableCell(
				fmt.Sprintf("%d ns", v.TotalCycleTime),
			).SetAlign(tview.AlignRight),
		)

		table.SetCell(
			r, 4,
			tview.NewTableCell(
				fmt.Sprintf("%d ns", v.TotalCycleTime/v.TotalTimesExecuted),
			).SetAlign(tview.AlignRight),
		)

	}

	if err := app.SetRoot(table, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return nil

}
