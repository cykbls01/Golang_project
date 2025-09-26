package main

import (
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
	"strconv"
)

func output(index int, sheet *xlsx.Sheet) {
	inputVal := map[string]interface{}{
		"ROSTemplateFormatVersion": "2015-09-01",
		"Resources":                make(map[string]interface{}),
	}
	columns := len(sheet.Rows[0].Cells)
	// 加入表格数据
	for i := 1 + index; i < 36+index && i < len(sheet.Rows); i++ {
		row := sheet.Rows[i]
		obj := map[string]interface{}{
			"AllocatePublicIP":   false,
			"DeletionProtection": false,
		}
		for j := 0; j < columns; j++ {
			key := sheet.Rows[0].Cells[j].String()
			value := row.Cells[j].String()
			key = splitKey(key)
			if key == "编号" {
				continue
			} else if key == "SystemDiskSize" {
				var intValue int
				fmt.Sscanf(value, "%d", &intValue)
				obj[key] = intValue
			} else if key == "StorageSetPartitionNumber" {
				var intValue int
				fmt.Sscanf(value, "%d", &intValue)
				obj[key] = intValue
			} else {
				obj[key] = value
			}
		}
		resourceKey := fmt.Sprintf("ECS%d", i-1)
		inputVal["Resources"].(map[string]interface{})[resourceKey] = map[string]interface{}{
			"Type":       "ALIYUN::ECS::Instance",
			"Properties": obj,
		}
	}

	// 将结果写入 JSON 文件
	file, err := os.Create("output/data" + strconv.Itoa(index/35+1) + ".json")
	if err != nil {
		fmt.Println("创建文件时出错:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(inputVal)
	if err != nil {
		fmt.Println("写入 JSON 文件时出错:", err)
		return
	}
}

func main() {
	// 读取 Excel 文件
	xlFile, err := xlsx.OpenFile("./ros.xlsx")
	if err != nil {
		fmt.Println("读取 Excel 文件时出错:", err)
		return
	}

	// 假设只有一个工作表
	sheet := xlFile.Sheets[0]
	rows := len(sheet.Rows)
	for index := 1; index < rows; index += 35 {
		output(index, sheet)
	}
}

// 辅助函数，用于分割键名
func splitKey(key string) string {
	for i, r := range key {
		if r == '(' {
			return key[:i]
		}
	}
	return key
}
