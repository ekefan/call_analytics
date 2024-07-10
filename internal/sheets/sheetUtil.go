package sheets

import (
	"fmt"

	excel "github.com/xuri/excelize/v2"
)

// SaveResp fields for creating a customer file error response
type SaveResp struct {
	ErrMsg     error
	SuccessMsg string
}

// errorResp custom error response for file operations
func errorResp(err error, msg string) SaveResp {
	return SaveResp{
		ErrMsg:     fmt.Errorf("%v: %s", msg, err),
		SuccessMsg: "",
	}
}

// AppendData gets the row number and adds input call entry data to that row.
func AppendData(f *excel.File, filePath string, callEntry CallEntryArgs) (string, error) {
	// Determine the next empty row
	sheetName := "Sheet1"
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return "", fmt.Errorf("failed to get rows: %v", err)
	}
	nextRow := len(rows) + 1

	// Add new data
	date := string(callEntry.DateOfEntry.Format("2006/01/02"))
	callDuration := string(callEntry.CallTime)
	supportPersonel := string(callEntry.SupportPersonel)
	incident := string(callEntry.Incident)
	resolution := string(callEntry.Resolution)
	comment := string(callEntry.Comment)

	cellValues := []string{date, callDuration, supportPersonel, incident, resolution, comment}

	err = addNewData(f, filePath, sheetName, nextRow, cellValues)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return "success", nil
}

// addNewData adds new call entry to the active worksheet
// it update the column size based on the length of the input data
func addNewData(f *excel.File, filePath, sheetName string, row int, cellValues []string) error {
	column := rune('A')
	var columnName string
	for _, value := range cellValues {
		//set cell value
		columnName = string(column)
		fmt.Println(columnName)
		cell := fmt.Sprintf("%s%d", columnName, row)
		f.SetCellValue(sheetName, cell, value)

		// get current width of the column to adjust if cellvalue can not be contained in it
		currentWidth, err := f.GetColWidth(sheetName, columnName)
		if err != nil {
			return err
		}
		// check the input value width
		inputWidth := float64(len(value))

		// update the current width if the input width it more than it
		if inputWidth > currentWidth {
			if err := f.SetColWidth(sheetName, columnName, columnName, inputWidth); err != nil {
				return err
			}
		}
		column++
	}

	// Save the file
	if err := f.SaveAs(filePath); err != nil {
		return err
	}
	return nil
}
