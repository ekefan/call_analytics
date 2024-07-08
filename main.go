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
	w.Resize(fyne.Size{Width: 800, Height: 650})

	// ==============widgets========================================
	// clk is the clock widget for counting the call duration
	clk := widget.NewLabel("00:00:00")

	// btns defined with out any callback
	startBtn := widget.NewButton("Start", nil)
	stopBtn := widget.NewButton("Stop", nil)
	saveBtn := widget.NewButton("Save", nil)

	// startBtn button used to start the clock
	startBtn.OnTapped = func() {
		// check if the timer is running
		if !timeRunning {
			timeRunning = true
			// entryContent.Hide() //
			calltime.StartTimer(clk)
			startBtn.Disable() // Disable the start button
			saveBtn.Enable()
			saveClicked = false
		}
	}

	// button used to trigger events that help to collect data
	stopBtn.OnTapped = func() {
		if timeRunning {
			timeRunning = false
			durationPlaceHolder := fmt.Sprintf("Call duration: %v", clk.Text)
			duration.SetText(durationPlaceHolder)
			entryContent.Show()
			calltime.StopTimer()
			timerContent.Hide()
		}
	}

	// saveBtn button used to save the call entry
	saveBtn.OnTapped = func() {
		// get text content from each widget...
		saveArgs := sheets.CallEntryArgs{ // Convert the call entry function to CallEntry datastructure
			CallTime:        clk.Text,
			SupportPersonel: supportPersonel.Text,
			Incident:        incident.Text,
			Resolution:      resolution.Text,
			Comment:         comment.Text,
			DateOfEntry:     time.Now(),
		}
		// pass the data into the saveCallEntry function
		if !saveClicked {
			// convert saveArgs to slice... check for empty place holder...
			//if empty... show... enter all fields widget
			success := sheets.SaveCallEntry(saveArgs)
			if success.ErrMsg != nil {
				fmt.Printf("Send Message to error widget")
			}
			entrySaved.Show()
			saveClicked = true
			clk.SetText("00:00:00")
			timerContent.Show()
			startBtn.Enable() // Enable the start button again
			time.Sleep(2 * time.Second)
			entrySaved.Hide()
			duration.SetText("")
			supportPersonel.SetText("")
			incident.SetText("")
			resolution.SetText("")
			comment.SetText("")
			entryContent.Hide()
		}
	}

	// widget for taking the call summary and details
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
