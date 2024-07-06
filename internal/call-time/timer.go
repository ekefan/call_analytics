package calltime

import (
	"time"

	"fyne.io/fyne/v2/widget"
	"fmt"
)


var (
	ticks *time.Ticker
	callEnded chan bool
)

func StartTimer(clk *widget.Label) {
	callEnded = make(chan bool)
	ticks = time.NewTicker(1 * time.Second)
	start := time.Now()

	go func(clk *widget.Label){
		for {
			select {
			case <-callEnded :
				return
			case t := <-ticks.C:
				elapsed := t.Sub(start)
				hours := int(elapsed.Hours())
				minutes := int(elapsed.Minutes()) % 60
				seconds := int(elapsed.Seconds()) % 60
				format := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
				clk.SetText(format)
			}

		}
	}(clk)


}

func StopTimer() {
	ticks.Stop()
	callEnded <- true
}
