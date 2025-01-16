package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func WriteSingleColumn(filePath string, sheetName string, strings []string) error {
	f := excelize.NewFile()

	for rowIndex, str := range strings {
		cellName := fmt.Sprintf("A%d", rowIndex+1)
		err := f.SetCellValue(sheetName, cellName, str)
		if err != nil {
			return err
		}
	}

	if err := f.SaveAs(filePath); err != nil {
		return err
	}
	return nil
}

func ReadRows(path, sheetName string) [][]string {
	f, _ := excelize.OpenFile(path)
	rows, _ := f.GetRows(sheetName)
	return rows
}

func WriteRows(path, sheetName string, rows [][]string) error {
	f := excelize.NewFile()
	for rowIndex, str := range rows {
		cellName := fmt.Sprintf("A%d", rowIndex+1)
		err := f.SetSheetRow(sheetName, cellName, &str)
		if err != nil {
			return err
		}
	}
	if err := f.SaveAs(path); err != nil {
		return err
	}
	return nil
}
