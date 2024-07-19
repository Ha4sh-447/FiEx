package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Ha4sh-447/FiEx/internal/ui"
)

func main() {
	a := app.New()
	w := a.NewWindow("FiEx")
	_, err := ui.InitScreen(w)

	if err != nil {
		fyne.LogError("Can't load window", err)
		return
	}

	w.Resize(fyne.NewSize(1200, 600))
	w.ShowAndRun()

}
