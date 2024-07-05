package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()

	w := myApp.NewWindow("Support Analytics")

	// Creating widgets for the timer, button, and data entry
	// ======widgets-======
	clk := widget.NewLabel("00:00:00")
	startBtn := widget.NewButton("Start", nil)
	stopBtn := widget.NewButton("stop", nil)
	duration := widget.NewLabel("Call duration: 00:00:00")
	duration.SetText(clk.Text)
	engineerDet := widget.NewEntry()
	engineerDet.PlaceHolder = "Support Engineer: "
	detail1 := widget.NewEntry()
	detail1.PlaceHolder = "Enter detail"

	timerContent := container.New(layout.NewVBoxLayout(), container.New(layout.NewCenterLayout(), clk),
	container.New(layout.NewAdaptiveGridLayout(2), startBtn, stopBtn))
	w.SetContent(timerContent)
	w.Show()
	myApp.Run()
	// timerContent
	// entryContent
}
