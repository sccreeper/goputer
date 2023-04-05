package main

import "fyne.io/fyne/v2"

func main() {

	Win.Resize(fyne.NewSize(640, 480))
	Win.SetMainMenu(MainMenu)
	Win.SetContent(CodeArea)
	Win.ShowAndRun()

}
