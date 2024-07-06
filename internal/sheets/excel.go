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
	SupportEngineer string
	CallDetail      string
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
		f.SetCellValue("Sheet1", "B1", "Call Duration")
		f.SetCellValue("Sheet1", "C1", "Support Engineer")
		f.SetCellValue("sheet1", "D1", "Details")

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


func errorResp(err error, msg string, ) SaveResp{
	return SaveResp{
		ErrMsg: fmt.Errorf("%v: %s", msg, err),
		SuccessMsg: "",
	}
}
func appendData(f *excel.File, filePath string, callEntry CallEntryArgs) (string, error){
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
	supportEngineer := callEntry.SupportEngineer
	details := callEntry.CallDetail

	fmt.Println(details, supportEngineer)
	// dateTime := time.Now().Format("2006-01-02 15:04:05")

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", nextRow), date)
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", nextRow), callDuration)
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", nextRow), supportEngineer)
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", nextRow), details)

	// Save the file
	if err := f.SaveAs(filePath); err != nil {
		
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return "success", nil
}
