package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FrontendInfo struct {
	Info struct {
		Name        string `toml:"name"`
		Description string `toml:"description"`
		Authour     string `toml:"authour"`
		Repository  string `toml:"repository"`
		IsPlugin    bool   `toml:"is_plugin"`
	} `toml:"info"`

	PluginSO string
}

var (
	App = app.New()
	Win = App.NewWindow("goputer Launcher")

	BtnOpenCode = widget.NewButton("Run Code", run_code)
	BtnRunCode  = widget.NewButton("Open Code", open_code)

	LblSelectedCode = widget.NewLabel(CodePath)

	LstFrontend = widget.NewList(

		func() int {
			return len(Frontends)
		},
		func() fyne.CanvasObject {
			return container.New(layout.NewVBoxLayout(),
				widget.NewLabel("Name"),
				widget.NewLabel("Description"),
			)
		},

		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[0].(*widget.Label).SetText(Frontends[i].Info.Name)
			o.(*fyne.Container).Objects[1].(*widget.Label).SetText(Frontends[i].Info.Description)
		},
	)
)

var (
	Code       []byte
	CodePath   string = "<no code selected>"
	CodeOpened bool   = false
)

var (
	Frontends        []FrontendInfo
	SelectedFrontend int
)
