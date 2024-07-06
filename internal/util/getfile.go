package util

import (
	"time"
	"os"
	"path/filepath"
	"fmt"
)
func GetFilePath(appName string) (string, error) {
	
	// Get the current month and year
	now := time.Now()
	month := now.Month().String()
	year := now.Year()

	// Get the path to the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get home directory: %s", err)
	}

	// Build the path to the folder in the home directory
	appFolder := filepath.Join(homeDir, appName)

	// Create the directory if it doesn't exist
	if _, err := os.Stat(appFolder); os.IsNotExist(err) {
		err := os.MkdirAll(appFolder, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %v", err)
		}
	}

	// Create the file path
	fileName := fmt.Sprintf("%s_%d.xlsx", month, year)
	filePath := filepath.Join(appFolder, fileName)
	return filePath, nil
}