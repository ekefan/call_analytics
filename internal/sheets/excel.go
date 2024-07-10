package sheets

import (
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

type Style struct {
	Font struct {
		Bold bool `json:"bold"`
	} `json:"font"`
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
			return errorResp(err, "can't append data to sheet")
		}
	} else {
		// Create a new file
		f = excel.NewFile()
		// Create a new sheet
		index, err := f.NewSheet("Sheet1")
		if err != nil {
			return errorResp(err, "failed to create file sheet")
		}

		// Create bold style for headers
		style := &excel.Style{
			Font: &excel.Font{
				Bold: true,
			},
		}
		styleID, err := f.NewStyle(style)
		if err != nil {
			return errorResp(err, "failed to create style")
		}

		// Set headers
		headers := []string{"A1", "B1", "C1", "D1", "E1", "F1"}
		columns := []string{"A", "B", "C", "D", "E", "F"}
		headerValues := []string{"Date", "Support Personel", "Call Duration", "Incident", "Resolution", "Comment"}

		for i, header := range headers {
			//set cell value
			f.SetCellValue("Sheet1", header, headerValues[i])

			// Apply style to header
			f.SetCellStyle("Sheet1", header, header, styleID)

			width, _ := f.GetColWidth("sheet1", columns[i])
			fmt.Printf("col width of %s: %v\n", header, width)
			err := f.SetColWidth("sheet1", columns[i], columns[i], 20.00)
			if err != nil {
				fmt.Printf("errrrrr: %v", err)
			}
			fmt.Printf("after changinging the width\n")
			width, _ = f.GetColWidth("sheet1", columns[i])
			fmt.Printf("col width of %s: %v\n", header, width)
		}

		f.SetActiveSheet(index)
		// Save the file
		if err := f.SaveAs(filePath); err != nil {
			return errorResp(err, "failed to create file:")
		}
		// Append data
		_, err = appendData(f, filePath, callEntry)
		if err != nil {
			return errorResp(err, "can't append data to sheet")
		}
	}

	// Adjust the width of the columns
	// if err := adjustColumnWidth(f, "Sheet1", []string{"A", "B", "C", "D", "E", "F"}); err != nil {
	// 	return errorResp(err, "failed to adjust column width")
	// }

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

// func adjustColumnWidth(f *excel.File, sheetName string, columns []string) error {
// 	// Set all columns to a larger width
// 	const largeWidth = 30 // Adjust as needed
// 	f.SetColWidth(sheetName, columns[0], columns[len(columns)-1], largeWidth)

// 	return nil
// }
