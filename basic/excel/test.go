package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
)

func main() {
	// 打开 Excel 文件
	filePath := "data.xlsx" // 替换为你的 Excel 文件路径
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatalf("Failed to open Excel file: %v", err)
	}

	// 查找名为 "CCE" 的工作表
	var sheetName = "科技园CCE"
	// 读取工作表中的所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("Failed to get rows from sheet '%s': %v", sheetName, err)
	}
	for _, row := range rows {
		fmt.Println(row)
	}
}
