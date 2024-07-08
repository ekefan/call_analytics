package sheets

import (
	// "encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/ekefan/call_analytics/internal/util"
	excel "github.com/xuri/excelize/v2"
)

const AppName string = "VATSCA"

type CallEntryArgs struct {
	CallTime        string
	SupportPersonel string
	Incident        string
	Resolution      string
	Comment         string
	DateOfEntry     time.Time
}
type SaveResp struct {
	ErrMsg     error
	SuccessMsg string
}

func SaveCallEntry(callEntry CallEntryArgs) SaveResp {
	// Get file path based on the date
	filePath, err := util.GetFilePath(AppName)
	if err != nil {
		return SaveResp{ErrMsg: err, SuccessMsg: ""}
	}

	fileExist := true
	// open the file and if one doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fileExist = false
	}

	var f *excel.File
	if fileExist {
		//open the existing File
		f, err = excel.OpenFile(filePath)
		if err != nil {

			return errorResp(err, "failed to open file:")
		}
		// Append data
		_, err := appendData(f, filePath, callEntry)
		if err != nil {
			return errorResp(err, "cant append data to sheet")
		}
	} else {
		// Create a new file
		f = excel.NewFile()
		// Create a new sheet
		index, err := f.NewSheet("Sheet1")
		if err != nil {
			return errorResp(err, "failed to create file sheet")
		}
		// Set headers
		f.SetCellValue("Sheet1", "A1", "Date")
		f.SetCellValue("Sheet1", "B1", "Support Personel")
		f.SetCellValue("Sheet1", "C1", "Call Duration")
		f.SetCellValue("sheet1", "D1", "Incident")
		f.SetCellValue("sheet1", "E1", "Resolution")
		f.SetCellValue("sheet1", "F1", "Comment")

		f.SetActiveSheet(index)
		// Save the file
		if err := f.SaveAs(filePath); err != nil {
			return errorResp(err, "failed to create file:")
		}
		// Append data
		_, err = appendData(f, filePath, callEntry)
		if err != nil {
			return errorResp(err, "cant append data to sheet")
		}
	}

	return SaveResp{ErrMsg: nil, SuccessMsg: "Saved successfully"}
}

func errorResp(err error, msg string) SaveResp {
	return SaveResp{
		ErrMsg:     fmt.Errorf("%v: %s", msg, err),
		SuccessMsg: "",
	}
}
func appendData(f *excel.File, filePath string, callEntry CallEntryArgs) (string, error) {
	// Determine the next empty row
	sheetName := "Sheet1"
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return "", fmt.Errorf("failed to get rows: %v", err)
	}
	nextRow := len(rows) + 1

	// Add new data
	date := callEntry.DateOfEntry
	callDuration := callEntry.CallTime
	supportPersonel := callEntry.SupportPersonel
	incident := callEntry.Incident
	resolution := callEntry.Resolution
	comment := callEntry.Comment

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", nextRow), date)
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", nextRow), supportPersonel)
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", nextRow), callDuration)
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", nextRow), incident)
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", nextRow), resolution)
	f.SetCellValue(sheetName, fmt.Sprintf("F%d", nextRow), comment)

	// Save the file
	if err := f.SaveAs(filePath); err != nil {

		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return "success", nil
}
