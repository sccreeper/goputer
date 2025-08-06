package main

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/profiler"
	"sccreeper/goputer/pkg/util"
	"slices"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"
)

const (
	sortModeAddress int = iota
	sortModeTimesExecuted
	sortModeTotalExecutionTime
	sortModeMeanExecutionTime
)

const menuTextString string = "[red]F1:[white] Sorting attribute [red]F2:[white] Sorting direction [red]F3:[white] Toggle conditional formatting [red]F4:[white] Toggle grouping"

var sortMode int
var sortAscending bool = true

var useConditionalFormatting bool = true

var useGrouping bool

var profileEntriesSlice []profiler.ProfileEntry
var profileEntriesSliceGrouped []profiler.ProfileEntry

var (
	mainTable *tview.Table
	flexRoot  *tview.Flex
)

func formatColour[T util.Number](val T, min T, max T) tcell.Color {
	return tcell.NewRGBColor(
		255,
		int32(util.Lerp(1.0-util.Normalise(val, min, max), 0, 200)),
		int32(util.Lerp(1.0-util.Normalise(val, min, max), 0, 200)),
	)
}

func renderDefaultTableView() {

	sortArrow := "▼"

	if sortAscending {
		sortArrow = "▲"
	}

	var colOffset int

	if useGrouping {
		colOffset = -1
	}

	if !useGrouping {

		mainTable.SetCell(
			0, 0,
			tview.NewTableCell(
				"Address",
			).SetBackgroundColor(tcell.ColorWhite).SetTextColor(tcell.ColorBlack),
		)
		if sortMode == sortModeAddress {
			mainTable.GetCell(0, 0).
				SetBackgroundColor(tcell.ColorBlack).
				SetTextColor(tcell.ColorWhite).
				Text += fmt.Sprintf(" %s", sortArrow)
		}

	}

	mainTable.SetCell(
		0, colOffset+1,
		tview.NewTableCell(
			"Instruction",
		).SetBackgroundColor(tcell.ColorWhite).SetTextColor(tcell.ColorBlack),
	)

	mainTable.SetCell(
		0, colOffset+2,
		tview.NewTableCell(
			"Times executed",
		).SetBackgroundColor(tcell.ColorWhite).SetTextColor(tcell.ColorBlack),
	)
	if sortMode == sortModeTimesExecuted {
		mainTable.GetCell(0, colOffset+2).
			SetBackgroundColor(tcell.ColorBlack).
			SetTextColor(tcell.ColorWhite).
			Text += fmt.Sprintf(" %s", sortArrow)
	}

	mainTable.SetCell(
		0, colOffset+3,
		tview.NewTableCell(
			"Total execution time",
		).SetBackgroundColor(tcell.ColorWhite).SetTextColor(tcell.ColorBlack),
	)
	if sortMode == sortModeTotalExecutionTime {
		mainTable.GetCell(0, colOffset+3).
			SetBackgroundColor(tcell.ColorBlack).
			SetTextColor(tcell.ColorWhite).
			Text += fmt.Sprintf(" %s", sortArrow)
	}

	mainTable.SetCell(
		0, colOffset+4,
		tview.NewTableCell(
			"Mean execution time",
		).SetBackgroundColor(tcell.ColorWhite).SetTextColor(tcell.ColorBlack),
	)
	if sortMode == sortModeMeanExecutionTime {
		mainTable.GetCell(0, colOffset+4).
			SetBackgroundColor(tcell.ColorBlack).
			SetTextColor(tcell.ColorWhite).
			Text += fmt.Sprintf(" %s", sortArrow)
	}

	setTableData()

}

func setTableData() {

	var sliceToUse []profiler.ProfileEntry

	if useGrouping {
		sliceToUse = profileEntriesSliceGrouped
	} else {
		sliceToUse = profileEntriesSlice
	}

	minTotalCycleTime := slices.MinFunc(sliceToUse, func(a profiler.ProfileEntry, b profiler.ProfileEntry) int {
		return cmp.Compare(a.TotalCycleTime, b.TotalCycleTime)
	}).TotalCycleTime
	maxTotalCycleTime := slices.MaxFunc(sliceToUse, func(a profiler.ProfileEntry, b profiler.ProfileEntry) int {
		return cmp.Compare(a.TotalCycleTime, b.TotalCycleTime)
	}).TotalCycleTime

	temp := slices.MinFunc(sliceToUse, func(a profiler.ProfileEntry, b profiler.ProfileEntry) int {
		return cmp.Compare(a.TotalCycleTime/a.TotalTimesExecuted, b.TotalCycleTime/b.TotalTimesExecuted)
	})
	minSingleCycleTime := temp.TotalCycleTime / temp.TotalTimesExecuted

	temp = slices.MaxFunc(sliceToUse, func(a profiler.ProfileEntry, b profiler.ProfileEntry) int {
		return cmp.Compare(a.TotalCycleTime/a.TotalTimesExecuted, b.TotalCycleTime/b.TotalTimesExecuted)
	})
	maxSingleCycleTime := temp.TotalCycleTime / temp.TotalTimesExecuted

	minTimesExecuted := slices.MinFunc(sliceToUse, func(a profiler.ProfileEntry, b profiler.ProfileEntry) int {
		return cmp.Compare(a.TotalTimesExecuted, b.TotalTimesExecuted)
	}).TotalTimesExecuted
	maxTimesExecuted := slices.MaxFunc(sliceToUse, func(a profiler.ProfileEntry, b profiler.ProfileEntry) int {
		return cmp.Compare(a.TotalTimesExecuted, b.TotalTimesExecuted)
	}).TotalTimesExecuted

	var colOffset int
	if useGrouping {
		colOffset = -1
	}

	for r, v := range sliceToUse {

		if !useGrouping {
			mainTable.SetCell(
				r+1, 0,
				tview.NewTableCell(
					fmt.Sprintf("0x%08X", v.Address),
				).SetTextColor(
					tcell.ColorGreen,
				).SetAlign(tview.AlignCenter),
			)
		}

		var itnString string

		if !useGrouping {
			itnString, _ = compiler.DecodeInstructionString(v.Instruction[:])
		} else {
			itnString = constants.InstructionIntsReversed[uint32(v.Instruction[0])]
		}

		mainTable.SetCell(
			r+1, colOffset+1,
			tview.NewTableCell(
				itnString,
			).SetMaxWidth(24),
		)

		mainTable.SetCell(
			r+1, colOffset+2,
			tview.NewTableCell(
				strconv.FormatInt(int64(v.TotalTimesExecuted), 10),
			).SetAlign(tview.AlignRight),
		)

		mainTable.SetCell(
			r+1, colOffset+3,
			tview.NewTableCell(
				fmt.Sprintf("%d ns", v.TotalCycleTime),
			).SetAlign(tview.AlignRight),
		)

		mainTable.SetCell(
			r+1, colOffset+4,
			tview.NewTableCell(
				fmt.Sprintf("%d ns", v.TotalCycleTime/v.TotalTimesExecuted),
			).SetAlign(tview.AlignRight),
		)

		if useConditionalFormatting {
			mainTable.GetCell(r+1, colOffset+2).SetTextColor(formatColour(v.TotalTimesExecuted, minTimesExecuted, maxTimesExecuted))
			mainTable.GetCell(r+1, colOffset+3).SetTextColor(formatColour(v.TotalCycleTime, minTotalCycleTime, maxTotalCycleTime))
			mainTable.GetCell(r+1, colOffset+4).SetTextColor(formatColour(v.TotalCycleTime/v.TotalTimesExecuted, minSingleCycleTime, maxSingleCycleTime))
		}

	}
}

func changeSortingOrder(changeSortMode bool) {

	if changeSortMode {
		sortMode++

		if sortMode > sortModeMeanExecutionTime {
			if useGrouping {
				sortMode = sortModeTimesExecuted
			} else {
				sortMode = sortModeAddress
			}
		}

		var sliceToUse []profiler.ProfileEntry
		if useGrouping {
			sliceToUse = profileEntriesSliceGrouped
		} else {
			sliceToUse = profileEntriesSlice
		}

		switch sortMode {
		case sortModeAddress:
			slices.SortFunc(sliceToUse, func(a profiler.ProfileEntry, b profiler.ProfileEntry) int {
				return cmp.Compare(a.Address, b.Address)
			})
		case sortModeTimesExecuted:
			slices.SortFunc(sliceToUse, func(a profiler.ProfileEntry, b profiler.ProfileEntry) int {
				return cmp.Compare(a.TotalTimesExecuted, b.TotalTimesExecuted)
			})
		case sortModeMeanExecutionTime:
			slices.SortFunc(sliceToUse, func(a profiler.ProfileEntry, b profiler.ProfileEntry) int {
				return cmp.Compare(a.TotalCycleTime/a.TotalTimesExecuted, b.TotalCycleTime/b.TotalTimesExecuted)
			})
		case sortModeTotalExecutionTime:
			slices.SortFunc(sliceToUse, func(a profiler.ProfileEntry, b profiler.ProfileEntry) int {
				return cmp.Compare(a.TotalCycleTime, b.TotalCycleTime)
			})
		}
	}

	if useGrouping {
		if sortAscending && !changeSortMode {
			slices.Reverse(profileEntriesSliceGrouped)
		} else if !sortAscending {
			slices.Reverse(profileEntriesSliceGrouped)
		}
	} else {
		if sortAscending && !changeSortMode {
			slices.Reverse(profileEntriesSlice)
		} else if !sortAscending {
			slices.Reverse(profileEntriesSlice)
		}
	}
}

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

	profileEntriesSlice = make([]profiler.ProfileEntry, 0, len(profileData.Data))

	for _, v := range profileData.Data {
		profileEntriesSlice = append(profileEntriesSlice, *v)
	}

	profileEntriesSliceGrouped = make([]profiler.ProfileEntry, 0, constants.HighestInstruction)

	instructionIndexes := make([]uint32, 0, constants.HighestInstruction)

	for _, v := range profileData.Data {

		if slices.Contains(instructionIndexes, uint32(v.Instruction[0])) {

			itnIndex := slices.IndexFunc(
				profileEntriesSliceGrouped,
				func(p profiler.ProfileEntry) bool {
					return p.Instruction[0] == v.Instruction[0]
				},
			)

			profileEntriesSliceGrouped[itnIndex].TotalCycleTime += v.TotalCycleTime
			profileEntriesSliceGrouped[itnIndex].TotalTimesExecuted += v.TotalTimesExecuted

		} else {

			instructionIndexes = append(instructionIndexes, uint32(v.Instruction[0]))
			profileEntriesSliceGrouped = append(profileEntriesSliceGrouped, *v)

		}

	}

	sortMode = -1
	changeSortingOrder(true)

	// Setup terminal app

	app := tview.NewApplication().SetTitle("goputer Profiler")

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyF1 {
			changeSortingOrder(true)
			renderDefaultTableView()

			return nil
		} else if event.Key() == tcell.KeyF2 {
			sortAscending = !sortAscending

			changeSortingOrder(false)
			renderDefaultTableView()

			return nil
		} else if event.Key() == tcell.KeyF3 {
			useConditionalFormatting = !useConditionalFormatting

			renderDefaultTableView()
		} else if event.Key() == tcell.KeyF4 {
			useGrouping = !useGrouping
			sortMode = sortModeTimesExecuted

			mainTable.Clear()

			changeSortingOrder(false)
			renderDefaultTableView()
		}

		return event
	})

	flexRoot = tview.NewFlex().SetDirection(tview.FlexRow)

	tableContainer := tview.NewBox().SetBorder(true).SetTitle("Instructions")

	mainTable = tview.NewTable().SetFixed(1, 5)
	mainTable.Box = tableContainer

	menuText := tview.NewTextView().SetText(menuTextString)
	menuText.SetBorder(true)
	menuText.SetDynamicColors(true)

	flexRoot.AddItem(mainTable, 0, 3, false)
	flexRoot.AddItem(menuText, 5, 1, false)

	renderDefaultTableView()

	if err := app.SetRoot(flexRoot, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return nil

}
