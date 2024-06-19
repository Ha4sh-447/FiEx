package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Ha4sh-447/fileExp/internal/ui"
)

func main() {
	a := app.New()
	w, err := ui.InitScreen(a)

	if err != nil {
		fyne.LogError("Can't load window", err)
		return
	}

	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()

}
