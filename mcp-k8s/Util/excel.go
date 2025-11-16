package Util

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
)

func Output(data interface{}, filename string) error {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("input must be a slice")
	}

	valType := val.Type().Elem()
	if valType.Kind() != reflect.Struct {
		return fmt.Errorf("slice element type must be a struct")
	}

	f := excelize.NewFile()
	index, _ := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)

	fields := make([]string, valType.NumField())
	for i := 0; i < valType.NumField(); i++ {
		fields[i] = valType.Field(i).Name
	}

	for i, field := range fields {
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string('A'+byte(i)), 1), field)
	}

	for i := 0; i < val.Len(); i++ {
		for j := 0; j < valType.NumField(); j++ {
			valueField := val.Index(i).Field(j)
			f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string('A'+byte(j)), i+2), valueField.Interface())
		}
	}

	return f.SaveAs(filename)
}
