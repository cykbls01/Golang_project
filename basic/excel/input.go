package excel

import (
	"github.com/xuri/excelize/v2"
	"log"
)

func GetStrings() []string {
	f, err := excelize.OpenFile("data.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Fatal(err)
	}
	ids := []string{}
	for _, row := range rows {
		if len(row) > 0 {
			regionId := row[0]
			ids = append(ids, regionId)
		}
	}
	return ids[1:]
}
