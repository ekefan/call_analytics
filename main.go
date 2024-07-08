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

// main, the entry point of the call support analytics program
func main() {

	// create instance of the fyne application
	myApp := app.New()

	var (
		entryContent    *fyne.Container
		timerContent    *fyne.Container
		duration        *widget.Label
		supportPersonel *widget.Entry
		incident        *widget.Entry
		resolution      *widget.Entry
		comment         *widget.Entry
		entrySaved      *widget.Label
		timeRunning     bool
		saveClicked     bool
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

	if timeRunning {
		startBtn.Disable()
	}
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
		saveArgs := sheets.CallEntryArgs{ // Convert the call entry function to CallEntry datastructure
			CallTime:        clk.Text,
			SupportPersonel: supportPersonel.Text,
			Incident:        incident.Text,
			Resolution:      resolution.Text,
			Comment:         comment.Text,
			DateOfEntry:     time.Now(),
		}
		//pass the data into the saveCallEntry function
		if !saveClicked {
			success := sheets.SaveCallEntry(saveArgs)
			if success.ErrMsg != nil {
				fmt.Printf("Send Message to error widget")
			}
			entrySaved.Show()
			saveClicked = true
			clk.SetText("00:00:00")
			timerContent.Show()
			time.Sleep(2 * time.Second)
			entrySaved.Hide()
			duration.SetText("")
			supportPersonel.SetText("")
			incident.SetText("")
			resolution.SetText("")
			comment.SetText("")
			entryContent.Hide()

		}
	})
	// wiget for taking the call summary and details
	duration = widget.NewLabel("Call duration: 00:00:00")
	entrySaved = widget.NewLabel("Saved Successfully")
	entrySaved.Hide()

	// the widget for entry of data
	supportPersonel = widget.NewEntry()
	supportPersonel.PlaceHolder = "Support Personel:"

	incident = widget.NewEntry()
	incident.PlaceHolder = "Incident:"

	resolution = widget.NewEntry()
	resolution.PlaceHolder = "Resolution:"

	comment = widget.NewEntry()
	comment.PlaceHolder = "Comment"

	//	===========================Containers=============================
	timerContent = container.NewBorder(
		nil, // top
		container.New(layout.NewAdaptiveGridLayout(2), startBtn, stopBtn), // bottom
		nil, // left
		nil, // right
		container.New(layout.NewCenterLayout(), clk), // center
	)
	entryContent = container.NewBorder(
		nil, // top
		container.New(layout.NewVBoxLayout(), entrySaved, saveBtn), // bottom
		nil, // left
		nil, // right
		container.New(layout.NewVBoxLayout(), duration, supportPersonel, incident, resolution, comment), // center
	)

	entryContent.Hide()
	w.SetContent(container.NewHSplit(entryContent, timerContent))
	w.Show()
	myApp.Run()
}
