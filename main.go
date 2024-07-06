package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	calltime "github.com/ekefan/call_analytics/internal/call-time"
	"github.com/ekefan/call_analytics/internal/sheets"
)

// timeRunning checks if the time is counting for an ongoing call
var (
	timeRunning bool
	saveClicked bool
)

// main, the entry point of the call support analytics program
func main() {

	// create instance of the fyne application
	myApp := app.New()

	var (
		entryContent *fyne.Container
		timerContent *fyne.Container
		duration     *widget.Label
		engineerDet  *widget.Entry
		detail1      *widget.Entry
	)

	// set a new window to display
	w := myApp.NewWindow("Support Analytics")

	// ==============widgets========================================
	// clk is the clock widget for counting the call duration
	clk := widget.NewLabel("00:00:00")

	// startBtn button used to start the clock
	startBtn := widget.NewButton("Start", func() {
		//check if the timer is running
		if !timeRunning {
			timeRunning = true
			entryContent.Hide()
			calltime.StartTimer(clk)
			saveClicked = false
		}
	})

	// button used to trigger eveents that helps to collect dataa
	stopBtn := widget.NewButton("stop", func() {
		if timeRunning {
			timeRunning = false
			durationPlaceHolder := fmt.Sprintf("Call duration: %v", clk.Text)
			duration.SetText(durationPlaceHolder)
			entryContent.Show()
			calltime.StopTimer()
			timerContent.Hide()
		}

	})
	saveBtn := widget.NewButton("save", func() {
		//get text content from each widget...
		saveArgs := sheets.CallEntryArgs{// Convert the call entry function to CallEntry datastructure
			CallTime:        clk.Text,
			SupportEngineer: engineerDet.Text,
			CallDetail:     detail1.Text,
			DateOfEntry: time.Now(),
		}
		//pass the data into the saveCallEntry function
		if !saveClicked {
			success := sheets.SaveCallEntry(saveArgs)
			if success.ErrMsg != nil {
				fmt.Printf("Send Message to error widget")
			}
			fmt.Printf("send Message to feedback widget")
			saveClicked = true
			clk.SetText("00:00:00")
			timerContent.Show()
			time.Sleep(2 * time.Second)
			duration.SetText("")
			engineerDet.SetText("")
			detail1.SetText("")
			entryContent.Hide()

		}
	})
	// wiget for taking the call summary and details
	duration = widget.NewLabel("Call duration: 00:00:00")

	// the widget for entry of data
	engineerDet = widget.NewEntry()
	engineerDet.PlaceHolder = "Support Engineer:"

	detail1 = widget.NewEntry()
	detail1.PlaceHolder = "Enter detail"

	//	===========================Containers=============================
	timerContent = container.NewBorder(
		nil, // top
		container.New(layout.NewAdaptiveGridLayout(2), startBtn, stopBtn), // bottom
		nil, // left
		nil, // right
		container.New(layout.NewCenterLayout(), clk), // center
	)
	entryContent = container.NewBorder(
		nil,     // top
		saveBtn, // bottom
		nil,     // left
		nil,     // right
		container.New(layout.NewVBoxLayout(), duration, engineerDet, detail1), // center
	)
	entryContent.Hide()
	w.SetContent(container.NewHSplit(entryContent, timerContent))
	w.Show()
	myApp.Run()
}
