package sheets

import (
	"fmt"
	"os"
	"time"

	"github.com/ekefan/call_analytics/internal/util"
	excel "github.com/xuri/excelize/v2"
)

const AppName string = "Call Panda"
const defaultColWidth float64 = 25.00

type CallEntryArgs struct {
	CallTime        string
	SupportPersonel string
	Incident        string
	Resolution      string
	Comment         string
	DateOfEntry     time.Time
}

func SaveCallEntry(callEntry CallEntryArgs) SaveResp {
	// Get file path based on the date
	filePath, err := util.GetFilePath(AppName)
	if err != nil {
		return SaveResp{ErrMsg: err, SuccessMsg: ""}
	}

	fileExist := true
	// check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fileExist = false
	}

	var f *excel.File
	if fileExist {
		//open the file if it exist
		f, err = excel.OpenFile(filePath)
		if err != nil {
			return errorResp(err, "failed to open file:")
		}
		// Append data
		_, err := AppendData(f, filePath, callEntry)
		if err != nil {
			return errorResp(err, "can't append data to sheet")
		}
	} else {
		// Or create a new file, create a sheet, style the headers, and append new data
		f = excel.NewFile()
		// Create a new sheet
		index, err := f.NewSheet("Sheet1")
		if err != nil {
			return errorResp(err, "failed to create file sheet")
		}

		// Create Style for headers
		style := &excel.Style{
			Font: &excel.Font{
				Bold: true,
				Size: 14,
			},
		}
		styleID, err := f.NewStyle(style)
		if err != nil {
			return errorResp(err, "failed to create style")
		}

		// Set headers
		headers := []string{"A1", "B1", "C1", "D1", "E1", "F1"}
		columns := []string{"A", "B", "C", "D", "E", "F"}
		headerValues := []string{"Date", "Call Duration", "Support Personel", "Incident", "Resolution", "Comment"}

		for i, header := range headers {
			//set cell value
			f.SetCellValue("Sheet1", header, headerValues[i])

			// Apply style to header
			f.SetCellStyle("Sheet1", header, header, styleID)

			// Increase default column width to prevent ####### in cells
			err := f.SetColWidth("sheet1", columns[i], columns[i], defaultColWidth)
			if err != nil {
				fmt.Printf("errrrrr: %v", err)
			}

		}

		f.SetActiveSheet(index)
		// Save the file
		if err := f.SaveAs(filePath); err != nil {
			return errorResp(err, "failed to create file:")
		}
		// Append data
		_, err = AppendData(f, filePath, callEntry)
		if err != nil {
			return errorResp(err, "can't append data to sheet")
		}
	}
	return SaveResp{ErrMsg: nil, SuccessMsg: "Saved successfully"}
}
