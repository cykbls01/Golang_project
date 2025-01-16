package excel

import (
	"fmt"
	"log"
	"strings"
)

func FilterRows(rows [][]string, columnName string, filterSubstring []string) ([][]string, error) {
	headers := rows[0]
	colIndex := -1
	for i, header := range headers {
		if strings.TrimSpace(strings.ToLower(header)) == strings.TrimSpace(strings.ToLower(columnName)) {
			colIndex = i
			break
		}
	}
	if colIndex == -1 {
		return rows, fmt.Errorf("column name %s not found", columnName)
	}

	// 过滤行
	var filteredRows [][]string
	filteredRows = append(filteredRows, rows[0])
	for _, row := range rows[1:] { // 跳过列名行
		cellValue := strings.TrimSpace(row[colIndex])
		flag := false
		for _, filter := range filterSubstring {
			if strings.Contains(strings.ToLower(cellValue), strings.ToLower(filter)) {
				flag = true
				break
			}
		}
		if flag {
			filteredRows = append(filteredRows, row)
		}
	}

	return filteredRows, nil
}

func FilterCols(rows [][]string, colName string) ([]string, error) {
	var result []string
	headerFound := false
	colIndex := -1
	for _, row := range rows {
		if !headerFound {
			for index, colData := range row {
				if colData == colName {
					headerFound = true
					colIndex = index
					break
				}
			}
			continue
		}
		if colIndex == -1 {
			break
		}
		result = append(result, row[colIndex])
	}

	return result, nil
}

func DiffRows(rows1, rows2 [][]string, col1, col2 string) [][]string {
	col1Index1, ok1 := findColumnIndexByName(rows1[0], col1)
	col2Index2, ok2 := findColumnIndexByName(rows2[0], col2)
	var filteredRows [][]string
	filteredRows = append(filteredRows, rows1[0])
	if !ok1 || !ok2 {
		log.Fatalf("Failed to find column 'col1' in file1 or 'col2' in file2")
	}

	for _, row1 := range rows1[1:] { // 跳过表头
		col1Value := row1[col1Index1]
		found := false
		for _, row2 := range rows2[1:] { // 跳过表头
			if row2[col2Index2] == col1Value {
				found = true
				break
			}
		}
		if !found {
			filteredRows = append(filteredRows, row1)
		}
	}
	return filteredRows
}

func findColumnIndexByName(header []string, columnName string) (int, bool) {
	for index, name := range header {
		if name == columnName {
			return index, true
		}
	}
	return -1, false
}
