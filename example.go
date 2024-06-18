package main

import (
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

// func Contai() *fyne.Container {

// 	green := color.NRGBA{R: 0, G: 180, B: 0, A: 255}

// 	text1 := canvas.NewText("Hello", green)
// 	text2 := canvas.NewText("Hellow", green)

// 	text2.Move(fyne.NewPos(20, 20))
// 	// content := container.NewWithoutLayout(text1, text2)
// 	content := container.New(layout.NewAdaptiveGridLayout(2), text1, text2)

// 	return content
// }

func Eg_canvas(w fyne.Window) {
	c := w.Canvas()

	blue := color.NRGBA{R: 0, G: 0, B: 180, A: 255}
	rect := canvas.NewRectangle(blue)

	c.SetContent(rect)

	go func() {
		time.Sleep(time.Second)
		green := color.NRGBA{R: 0, G: 180, B: 0, A: 255}
		rect.FillColor = green
		rect.Resize(fyne.NewSize(100, 100))
		rect.Move(fyne.NewPos(20, 20))
		rect.Refresh()
	}()
}

// content := contai()

// myWindow.SetContent(content)
// eg_canvas(myWindow)

// ----------------  SHOW SYSTEM TRAY----------
// if desk, ok := myApp.(desktop.App); ok {
// 	m := fyne.NewMenu("MyApp",
// 		fyne.NewMenuItem("Show", func() {
// 			myWindow.Show()
// 		}))
// 	desk.SetSystemTrayMenu(m)
// }

// myWindow.SetContent(widget.NewLabel("Fyne System Tray"))
// myWindow.SetCloseIntercept(func() {
// 	myWindow.Hide()
// })
// -------------------------------

// ---------------- container.AppTabs cookup---------------
// First tab content
// clock := widget.NewLabel("")
// updateTime(clock)
// go func() {
// 	for range time.Tick(time.Second) {
// 		updateTime(clock)
// 	}
// }()
// clockTabContent := container.NewVBox(
// 	widget.NewLabel("Current Time:"),
// 	clock,
// )

// // Second tab content
// greetingLabel := widget.NewLabel("Hello, World!")
// greetingButton := widget.NewButton("Click Me", func() {
// 	greetingLabel.SetText("Hello, Fyne!")
// })
// greetingTabContent := container.NewVBox(
// 	greetingLabel,
// 	greetingButton,
// )

// // Third tab content
// data := []string{"Item 1", "Item 2", "Item 3"}
// list := widget.NewList(
// 	func() int {
// 		return len(data)
// 	},
// 	func() fyne.CanvasObject {
// 		return widget.NewLabel("")
// 	},
// 	func(i widget.ListItemID, o fyne.CanvasObject) {
// 		o.(*widget.Label).SetText(data[i])
// 	},
// )
// listTabContent := container.NewVBox(
// 	widget.NewLabel("Data List:"),
// 	list,
// )

// // Create tabs
// tabs := container.NewAppTabs(
// 	container.NewTabItem("Clock", clockTabContent),
// 	container.NewTabItem("Greeting", greetingTabContent),
// 	container.NewTabItem("List", listTabContent),
// )
// tabs.SetTabLocation(container.TabLocationTop)

// // Set the content of the window to the AppTabs container
// myWindow.SetContent(tabs)
// ---------------------------------
