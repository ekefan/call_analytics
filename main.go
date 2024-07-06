package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	calltime "github.com/ekefan/call_analytics/internal/call-time"
)

// timeRunning checks if the time is counting for an ongoing call
var timeRunning bool

// main, the entry point of the call support analytics program
func main() {

	// create instance of the fyne application
	myApp := app.New()

	
	var entryContent *fyne.Container
	var duration *widget.Label

	// set a new window to display
	w := myApp.NewWindow("Support Analytics")

	// ======widgets-======
	// clk is the clock widget for counting the call duration
	clk := widget.NewLabel("00:00:00")

	// startBtn button used to start the clock
	startBtn := widget.NewButton("Start", func() {
		//check if the timer is running
		if !timeRunning {
			timeRunning = true
			entryContent.Hide()
			calltime.StartTimer(clk)

		}
	})

	// button used to trigger eveents that helps to collect dataa
	stopBtn := widget.NewButton("stop", func() {
		if timeRunning {
			timeRunning = false
			duration.SetText(clk.Text)
			clk.SetText("00:00:00")
			entryContent.Show()
			calltime.StopTimer()
		}
		
	})

	// wiget for taking the call summary and details
	duration = widget.NewLabel("Call duration: 00:00:00")
	duration.SetText(clk.Text)

	// the widget for entry of data
	engineerDet := widget.NewEntry()
	engineerDet.PlaceHolder = "Support Engineer:"

	detail1 := widget.NewEntry()
	detail1.PlaceHolder = "Enter detail"

	timerContent := container.NewVSplit(
	container.New(layout.NewCenterLayout(), clk), 
	container.New(layout.NewAdaptiveGridLayout(2), startBtn, stopBtn))

	entryContent = container.New(layout.NewVBoxLayout(), duration, engineerDet, detail1)
	entryContent.Hide()
	w.SetContent(container.NewHSplit(entryContent, timerContent))
	w.Show()
	myApp.Run()
}



