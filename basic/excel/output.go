package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func NewFile(name, column string, ids []string) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	index, _ := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)

	columnTitleRow := 1
	cellID := "A" + fmt.Sprintf("%d", columnTitleRow)
	f.SetCellValue("Sheet1", cellID, column)

	for idx, id := range ids {
		rowIndex := idx + 2
		cell := "A" + fmt.Sprintf("%d", rowIndex)
		f.SetCellValue("Sheet1", cell, id)
	}

	if err := f.SaveAs(name); err != nil {
		fmt.Println(err)
	}
}
