package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Ha4sh-447/fileExp/internal/ui"
)

func main() {
	a := app.New()
	w := a.NewWindow("File Explorer")
	w_new, err := ui.InitScreen(w)

	if err != nil {
		fyne.LogError("Can't load window", err)
		return
	}

	w_new.Resize(fyne.NewSize(800, 600))
	w_new.ShowAndRun()

}
