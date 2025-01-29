package main

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {

	//Inits the frontend list
	getFrontends()

	//Inits the UI
	LstFrontend.OnSelected = func(id widget.ListItemID) { SelectedFrontend = id }

	Win.SetContent(
		container.New(layout.NewVBoxLayout(),
			LstFrontend,
			LblSelectedCode,
			container.New(layout.NewGridLayout(2),
				BtnRunCode,
				BtnOpenCode,
			),
		),
	)

	Win.ShowAndRun()

}
