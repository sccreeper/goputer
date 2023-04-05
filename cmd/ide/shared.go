package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

var (
	App = app.New()
	Win = App.NewWindow("goputer IDE")

	MainMenu *fyne.MainMenu
	FileMenu *fyne.Menu

	CodeArea = widget.NewMultiLineEntry()
)

func init() {

	FileMenu = fyne.NewMenu("File", fyne.NewMenuItem("Save", func() {}), fyne.NewMenuItem("New", func() {}))
	MainMenu = fyne.NewMainMenu(FileMenu)

}
