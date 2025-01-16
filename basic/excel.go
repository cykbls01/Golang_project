package main

import (
	"basic/util/excel"
	"fmt"
)

func main() {
	filePath1 := "util/excel/imcloud.xlsx"
	filePath2 := "util/excel/4a.xlsx"
	sheetName := "sheet1"

	rows1 := excel.ReadRows(filePath1, sheetName)
	rows1, _ = excel.FilterRows(rows1, "归属产品", []string{"容器管理", "应用托管平台"})
	rows2 := excel.ReadRows(filePath2, sheetName)
	diff := excel.DiffRows(rows1, rows2, "主IP", "IP地址")
	fmt.Println(excel.WriteRows("util/excel/output.xlsx", "Sheet1", diff))
}
