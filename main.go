package main

import (
	"fmt"
	"time"
)

func main() {
	var duration time.Duration
	start := time.Now()

	// Channel to communicate commands
	cmd := make(chan string)

	// Goroutine to read user input
	go func() {
		var command string
		for {
			fmt.Print("cmd: ")
			fmt.Scanln(&command)
			cmd <- command
			if command == "stop" {
				break
			}
		}
	}()

	// Goroutine to update and print the timer
	go func() {
		for {
			time.Sleep(1 * time.Second)
			duration = time.Since(start)
			hours := int(duration.Hours())
			minutes := int(duration.Minutes()) % 60
			seconds := int(duration.Seconds()) % 60

			// Move cursor up one line, clear the line, and print the timer
			fmt.Printf("\033[1A\033[2K%02d : %02d : %02d\ncmd: ", hours, minutes, seconds)
		}
	}()

	// Wait for the "stop" command
	for {
		command := <-cmd
		if command == "stop" {
			break
		}
	}

	duration = time.Since(start) // Ensure duration is up to date
	fmt.Println("\nFinished")
	fmt.Printf("The entire call duration is: %v\n", duration)
}
